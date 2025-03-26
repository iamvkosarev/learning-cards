package main

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/app/usecase"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/auth"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	sso_pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load")
	}

	cfg := config.MustLoad()
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("error setting up logger: %v\n", err)
	}

	ctx := context.Background()

	dns := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVICE_NAME"),
		os.Getenv("DB_PORT_INTERNAL"), os.Getenv("DB_NAME"),
	)
	pool, err := postgres.NewPostgresPool(ctx, dns)
	if err != nil {
		log.Fatalf("error setting up postgres: %v\ndns: %v", err, dns)
	}
	groupRepository := sqlRepository.NewGroupRepository(pool)

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ssoConn, err := grpc.NewClient(cfg.SSO.HostAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	ssoClient := sso_pb.NewSSOClient(ssoConn)
	authService := auth.NewGRPCService(ssoClient)

	grpcServer := grpc.NewServer()
	groupUseCase := usecase.NewGroupUseCase(groupRepository, authService)
	learningCardsServer := server.NewServer(groupUseCase, logger)

	pb.RegisterLearningCardsServer(grpcServer, learningCardsServer)

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	err = pb.RegisterLearningCardsHandlerFromEndpoint(
		ctx, mux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting REST gateway on", cfg.Server.RESTPort)
	if err := http.ListenAndServe(cfg.Server.RESTPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
