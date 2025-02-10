package handlers

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "terraform-executor/api/proto"
	"terraform-executor/internal/executor"
	"terraform-executor/internal/health"
)

// StartServer starts the gRPC server.
func StartServer(ctx context.Context, address string) error {
	// Create a TCP listener on the specified address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and initialize the Executor service with context
	executorService, err := executor.NewExecutorService(ctx)
	if err != nil {
		return err
	}

	// Register the Executor service with the gRPC server
	pb.RegisterExecutorServer(grpcServer, executorService)

	// Create and register health service with dependencies
	healthService := health.NewHealthService(executorService.K8sClient, executorService.AWSClient)
	pb.RegisterHealthServer(grpcServer, healthService)

	// Enable gRPC reflection for easier client interaction
	reflection.Register(grpcServer)

	// Handle server shutdown when context is cancelled
	go func() {
		<-ctx.Done()
		log.Println("Initiating graceful shutdown of gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")
	}()

	// Start serving gRPC requests
	return grpcServer.Serve(listener)
}
