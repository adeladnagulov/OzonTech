package main

import (
	"hw_2/internal/product"
	"log"
	"net"

	pb "hw_2/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	file string = "data.json"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	reflection.Register(server)

	srv, err := product.NewProductService(file)
	if err != nil {
		log.Fatalf("failed to init service: %v", err)
	}
	pb.RegisterProductServer(server, srv)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
