package main

import (
	"fmt"
	"log"
	"terraform-executor/api/handlers"
)

func main() {
	address := ":50051"
	fmt.Printf("Starting gRPC server on %s\n", address)

	err := handlers.StartServer(address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
