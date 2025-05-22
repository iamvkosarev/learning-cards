package app

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context)
}

type LoadConfigFunc[TConfig any] func() (TConfig, error)
type NewServerFunc[TConfig any, TServer Server] func(ctx context.Context, cfg TConfig) (TServer, error)

func PrepareAndRunApp[TConfig any, TServer Server](
	newServer NewServerFunc[TConfig, TServer],
	loadConfig LoadConfigFunc[TConfig],
) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found or failed to load: %v", err)
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server, err := newServer(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	go func() {
		if err = server.Start(); err != nil {
			log.Printf("Failed to initialize server: %v", err)
			shutdown <- syscall.SIGTERM
		}
	}()

	<-shutdown
	log.Println("Shutting down gracefully")

	server.Shutdown(ctx)
	log.Println("Server stopped")
}
