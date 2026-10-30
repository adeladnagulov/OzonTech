package get

import (
	"context"
	"hw_2/internal/cart/domain/cart"
	product "hw_2/pkg/api"

	"google.golang.org/grpc"
)

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test

type cartGetter interface {
	Get(userId int64) (*cart.Cart, error)
}

type cartProductClient interface {
	GetProduct(ctx context.Context, in *product.GetProductRequest, opts ...grpc.CallOption) (*product.GetProductResponse, error)
}
