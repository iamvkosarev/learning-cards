package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Start() error
	Shutdown(ctx context.Context)
}

func RunServer(
	ctx context.Context,
	server Server,
) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Failed to initialize server: %v", err)
			shutdown <- syscall.SIGTERM
		}
	}()

	<-shutdown
	log.Println("Shutting down gracefully")

	server.Shutdown(ctx)
	log.Println("Server stopped")
}
