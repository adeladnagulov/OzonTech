package product

import (
	"context"
	pb "hw_2/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ProductService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if err := s.validate.Struct(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	product, ok := s.products[in.Sku]
	if !ok {
		return nil, status.Error(codes.NotFound, "product not found")
	}
	return &pb.GetProductResponse{Sku: in.Sku, Name: product.Name, Price: int32(product.Price)}, nil
}
