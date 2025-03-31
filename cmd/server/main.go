package main

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/app/usecase"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/auth"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	sso_pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

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

	// === Repositories connection ===
	dns := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVICE_NAME"),
		os.Getenv("DB_PORT_INTERNAL"), os.Getenv("DB_NAME"),
	)
	pool, err := postgres.NewPostgresPool(ctx, dns)
	if err != nil {
		log.Fatalf("error setting up postgres: %v", err)
	}
	groupRepo := sqlRepository.NewGroupRepository(pool)
	cardRepo := sqlRepository.NewCardRepository(pool)

	// === Auth connection ===
	var authService contracts.AuthVerifier
	authService = getAuthVerifier(err, cfg)

	// === UseCases ===
	groupUseCase := usecase.NewGroupUseCase(
		usecase.GroupUseCaseDeps{
			GroupReader:  groupRepo,
			GroupWriter:  groupRepo,
			AuthVerifier: authService,
		},
	)
	cardsUseCase := usecase.NewCardsUseCase(
		usecase.CardsUseCaseDeps{
			GroupReader:  groupRepo,
			CardWriter:   cardRepo,
			CardReader:   cardRepo,
			AuthVerifier: authService,
		},
	)

	// === Server ===
	learningCardsServer := server.NewServer(groupUseCase, cardsUseCase, logger)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.RecoveryInterceptor(logger)),
	)

	pb.RegisterLearningCardsServer(grpcServer, learningCardsServer)

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting gRPC server on", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	httpMux := http.NewServeMux()

	gwMux := runtime.NewServeMux()

	err = pb.RegisterLearningCardsHandlerFromEndpoint(
		ctx, gwMux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	const firstVersion = "/v1/"

	httpMux.Handle(firstVersion, gwMux)

	httpMux.HandleFunc(
		cfg.Server.RestPrefix+firstVersion, func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, cfg.Server.RestPrefix)
			r2 := new(http.Request)
			*r2 = *r
			r2.URL.Path = path
			gwMux.ServeHTTP(w, r2)
		},
	)

	corsHandler := middleware.CorsWithOptions(httpMux, cfg.Server.CorsOptions)

	httpAddr := fmt.Sprintf("0.0.0.0%s", cfg.Server.RESTPort)

	log.Println("Starting REST gateway on", httpAddr)
	if err := http.ListenAndServe(httpAddr, corsHandler); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

func getAuthVerifier(err error, cfg *config.Config) contracts.AuthVerifier {
	var authService contracts.AuthVerifier
	if cfg.SSO.UseLocal {
		authService = auth.NewLocalService(cfg.SSO.LocalUserId)
		return authService
	}

	ssoConn, err := grpc.NewClient(cfg.SSO.HostAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	ssoClient := sso_pb.NewSSOClient(ssoConn)
	authService = auth.NewGRPCService(ssoClient)
	return authService
}
