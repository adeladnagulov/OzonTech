package get_test

import (
	"context"
	"hw_2/internal/cart/domain/cart"
	"hw_2/internal/cart/usecase/get"
	pb "hw_2/pkg/api"
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetCart(t *testing.T) {
	t.Parallel()
	userId := int64(2)

	tests := []struct {
		name           string
		setup          func(repo *MockcartGetter, productClient *MockcartProductClient)
		wantСart       cart.Cart
		wantErr        error
		wantItemsCount int
	}{
		{
			name: "success and sort",
			setup: func(repo *MockcartGetter, productClient *MockcartProductClient) {
				repo.EXPECT().
					Get(gomock.Any()).
					Return(&cart.Cart{
						Items: []cart.Item{
							{Sku: 1, Name: "1", Count: 1, Price: 100},
							{Sku: 3, Name: "3", Count: 1, Price: 100},
							{Sku: 2, Name: "2", Count: 2, Price: 100},
						},
						TotalPrice: 400,
					}, nil)
				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.AssignableToTypeOf(&pb.GetProductRequest{}), gomock.Any()).
					DoAndReturn(func(ctx context.Context, req *pb.GetProductRequest, opts ...interface{}) (*pb.GetProductResponse, error) {
						switch req.Sku {
						case 1:
							return &pb.GetProductResponse{Sku: 1, Name: "1", Price: 100}, nil
						case 2:
							return &pb.GetProductResponse{Sku: 2, Name: "2", Price: 100}, nil
						case 3:
							return &pb.GetProductResponse{Sku: 3, Name: "3", Price: 100}, nil
						default:
							return nil, status.Error(codes.NotFound, "not found")
						}
					}).AnyTimes()
			},
			wantСart: cart.Cart{
				Items: []cart.Item{
					{Sku: 1, Name: "1", Count: 1, Price: 100},
					{Sku: 2, Name: "2", Count: 2, Price: 100},
					{Sku: 3, Name: "3", Count: 1, Price: 100},
				},
				TotalPrice: 400,
			},
			wantErr:        nil,
			wantItemsCount: 3,
		},
		{
			name: "cart not found",
			setup: func(repo *MockcartGetter, productClient *MockcartProductClient) {
				repo.EXPECT().
					Get(gomock.Any()).
					Return(nil, cart.ErrProductNotFound)
			},
			wantErr:        cart.ErrProductNotFound,
			wantItemsCount: 0,
		},
		{
			name: "success empty cart no product calls",
			setup: func(repo *MockcartGetter, productClient *MockcartProductClient) {
				repo.EXPECT().Get(userId).Return(&cart.Cart{
					Items:      []cart.Item{},
					TotalPrice: 0,
				}, nil)
			},
			wantСart:       cart.Cart{Items: []cart.Item{}, TotalPrice: 0},
			wantErr:        nil,
			wantItemsCount: 0,
		},
		{
			name: "error_product_not_found",
			setup: func(repo *MockcartGetter, productClient *MockcartProductClient) {
				repo.EXPECT().Get(userId).Return(&cart.Cart{
					Items: []cart.Item{{Sku: 999, Count: 1}},
				}, nil)

				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.AssignableToTypeOf(&pb.GetProductRequest{}), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			wantErr:        cart.ErrProductNotFound,
			wantItemsCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcartGetter(ctrl)
			mockProductClient := NewMockcartProductClient(ctrl)
			tt.setup(mockRepo, mockProductClient)
			uc := get.NewGetUseCase(mockRepo, mockProductClient)

			c, err := uc.GetCart(context.Background(), userId)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantСart, c)
				require.Len(t, c.Items, tt.wantItemsCount)
			}
		})
	}
}
