package main

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/handlers/cards"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcPort := ":50051"
	restPort := ":8080"

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cardHandler := cards.NewCardsHandler()

	pb.RegisterLearningCardsServer(grpcServer, cardHandler)

	go func() {
		log.Println("Starting gRPC server on", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	err = pb.RegisterLearningCardsHandlerFromEndpoint(
		ctx, mux, grpcPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("Starting REST gateway on", restPort)
	if err := http.ListenAndServe(restPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
