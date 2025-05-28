package main

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
	"github.com/joho/godotenv"
	"log"
)

const reviewsConfigPathEnvKey = "CONFIG_PATH"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load: %v", err)
	}

	cfg, err := config.Load[config.ReviewsConfig](reviewsConfigPathEnvKey)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := server.BuildReviewsServer(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	app.RunServer(ctx, server)
}
