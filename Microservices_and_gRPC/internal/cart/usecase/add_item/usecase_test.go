package additem_test

import (
	"context"
	"hw_2/internal/cart/domain/cart"
	additem "hw_2/internal/cart/usecase/add_item"
	pb "hw_2/pkg/api"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddItem(t *testing.T) {
	t.Parallel()
	userId, skuId := int64(1), int64(2)

	tests := []struct {
		name           string
		setup          func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart)
		req            cart.AddItemReq
		wantTotalPrice uint32
		wantItemsCount int
		wantErr        error
	}{
		{
			name: "success_add_new_item_in_empmy_cart",
			setup: func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart) {
				repo.EXPECT().
					Create(gomock.Any()).
					Return(testCart)
				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&pb.GetProductResponse{
						Sku:   2,
						Name:  "Test product",
						Price: 100,
					}, nil)
			},
			req:            cart.AddItemReq{Count: 2},
			wantTotalPrice: 200,
			wantItemsCount: 1,
			wantErr:        nil,
		},
		{
			name:    "failed validate req",
			setup:   func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart) {},
			req:     cart.AddItemReq{Count: 0},
			wantErr: cart.ErrBadRequest,
		},
		{
			name: "product not found",
			setup: func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart) {
				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			req:     cart.AddItemReq{Count: 2},
			wantErr: cart.ErrProductNotFound,
		},
		{
			name: "error grpc unavailable",
			setup: func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart) {
				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Unavailable, "service down"))
			},
			req:     cart.AddItemReq{Count: 2},
			wantErr: cart.ErrGRPCUnavailable,
		},
		{
			name: "success append item in cart",
			setup: func(repo *MockcartCreator, productClient *MockcartProductClient, testCart *cart.Cart) {
				testCart.Items = append(testCart.Items, cart.Item{
					Sku:   5,
					Name:  "newItem",
					Count: 3,
					Price: 150,
				})
				testCart.Items = append(testCart.Items, cart.Item{
					Sku:   2,
					Name:  "Test product",
					Count: 1,
					Price: 100,
				})
				testCart.TotalPrice = 550
				repo.EXPECT().
					Create(gomock.Any()).
					Return(testCart)
				productClient.EXPECT().
					GetProduct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&pb.GetProductResponse{
						Sku:   2,
						Name:  "Test product",
						Price: 100,
					}, nil)
			},
			req:            cart.AddItemReq{Count: 1},
			wantTotalPrice: 650,
			wantItemsCount: 2,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ctx := context.Background()

			mockRepo := NewMockcartCreator(ctrl)
			mockProductClient := NewMockcartProductClient(ctrl)
			testCart := &cart.Cart{
				Items:      []cart.Item{},
				TotalPrice: 0,
			}
			tt.setup(mockRepo, mockProductClient, testCart)
			uc := additem.NewAddUseCase(mockRepo, mockProductClient)

			err := uc.AddItem(tt.req, ctx, skuId, userId)
			if err != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Equal(t, tt.wantTotalPrice, testCart.TotalPrice)
				assert.Len(t, testCart.Items, tt.wantItemsCount)
			}
		})
	}
}
