package main

import (
	"hw_2/internal/cart"
	pb "hw_2/pkg/api"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	srv := cart.NewCartService(pb.NewProductClient(conn))
	handler := cart.NewHandler(srv)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddItem)
	mux.HandleFunc("GET /user/{user_id}/cart", srv.GetCart)

	log.Println("Starting cart client on :8080")
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
