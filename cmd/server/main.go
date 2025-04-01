package main

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	newApp, err := app.NewApp(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	go func() {
		if err := newApp.Run(); err != nil {
			log.Printf("Failed to initialize application: %v", err)
			shutdown <- syscall.SIGTERM
		}
	}()

	<-shutdown
	log.Println("Shutting down gracefully")

	newApp.Shutdown(ctx)
	log.Println("Server stopped")

}
