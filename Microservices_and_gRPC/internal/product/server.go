package product

import (
	"context"
	pb "hw_2/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.ProductServer
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) GetProduct(context.Context, *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
