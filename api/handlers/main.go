package handlers

import (
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "terraform-executor/api/proto"
	"terraform-executor/internal/executor"
)

// StartServer starts the gRPC server.
func StartServer(address string) error {
	// Create a TCP listener on the specified address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the Executor service with the gRPC server
	executorService := &executor.ExecutorService{}
	pb.RegisterExecutorServer(grpcServer, executorService)

	// Enable gRPC reflection for easier client interaction (e.g., grpcurl)
	reflection.Register(grpcServer)

	// Start serving gRPC requests
	return grpcServer.Serve(listener)
}
