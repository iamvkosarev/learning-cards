package main

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/joho/godotenv"
	"log"
)

const cardsConfigPathEnvKey = "CONFIG_PATH"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load: %v", err)
	}

	cfg, err := config.Load[config.CardsConfig](cardsConfigPathEnvKey)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cardsApp, err := app.NewCardsApp(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}
	cardsApp.Run(ctx)
}
