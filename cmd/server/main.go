package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"terraform-executor/api/handlers"
)

func main() {
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, initiating shutdown", sig)
		cancel()
	}()

	address := os.Getenv("LISTEN")
	if address == "" {
		address = ":50051"
	}
	fmt.Printf("Starting gRPC server on %s\n", address)

	if err := handlers.StartServer(ctx, address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
