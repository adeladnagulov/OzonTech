package main

import (
	additem "hw_2/internal/cart/http/cart/add_item"
	"hw_2/internal/cart/http/cart/get"
	addUseCase "hw_2/internal/cart/usecase/add_item"
	getUseCase "hw_2/internal/cart/usecase/get"
	repository "hw_2/internal/repository/cart"
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

	repo := repository.NewCartsRepo()
	addUc := addUseCase.NewAddUseCase(repo, pb.NewProductClient(conn))
	getUc := getUseCase.NewGetUseCase(repo, pb.NewProductClient(conn))

	addH := additem.NewHandler(addUc)
	getH := get.NewHandler(getUc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", addH.AddItem)
	mux.HandleFunc("GET /user/{user_id}/cart", getH.GetCart)

	log.Println("Starting cart client on :8080")
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
