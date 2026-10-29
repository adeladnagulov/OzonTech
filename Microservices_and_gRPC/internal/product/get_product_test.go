package product

import (
	"context"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "hw_2/pkg/api"
)

func TestProductService_GetProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		req          *pb.GetProductRequest
		initialMap   map[int64]*Product
		wantCode     codes.Code
		wantResponse *pb.GetProductResponse
	}{
		{
			name: "success product found",
			req:  &pb.GetProductRequest{Sku: 1},
			initialMap: map[int64]*Product{
				1: {Name: "Apple", Price: 100},
			},
			wantCode: codes.OK,
			wantResponse: &pb.GetProductResponse{
				Sku:   1,
				Name:  "Apple",
				Price: 100,
			},
		},
		{
			name: "error product not found",
			req:  &pb.GetProductRequest{Sku: 99},
			initialMap: map[int64]*Product{
				1: {Name: "Apple", Price: 100},
			},
			wantCode:     codes.NotFound,
			wantResponse: nil,
		},
		{
			name:         "error validation failed",
			req:          &pb.GetProductRequest{Sku: 0},
			initialMap:   map[int64]*Product{},
			wantCode:     codes.InvalidArgument,
			wantResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := &ProductService{
				products: tt.initialMap,
				validate: validator.New(),
			}

			resp, err := svc.GetProduct(context.Background(), tt.req)

			require.Equal(t, tt.wantCode, status.Code(err))

			if tt.wantCode == codes.OK {
				require.NotNil(t, resp)
				require.Equal(t, tt.wantResponse, resp)
			} else {
				require.Nil(t, resp)
			}
		})
	}
}
