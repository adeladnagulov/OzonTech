package product

import (
	"context"
	"encoding/json"
	"fmt"
	pb "hw_2/pkg/api"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.ProductServer

	products map[int64]*Product
}

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func NewProductService(file string) (*ProductService, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read products file: %w", err)
	}
	var products map[int64]*Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products JSON: %w", err)
	}

	return &ProductService{products: products}, nil
}

func (s *ProductService) GetProduct(context.Context, *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
