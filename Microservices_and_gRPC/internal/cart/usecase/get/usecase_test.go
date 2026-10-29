package get_test

import (
	"hw_2/internal/cart/domain/cart"
	"hw_2/internal/cart/usecase/get"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetCart(t *testing.T) {
	t.Parallel()
	userUd := int64(2)

	tests := []struct {
		name     string
		setup    func(repo MockcartGetter)
		wantСart cart.Cart
		wantErr  error
	}{
		{
			name: "success and sort",
			setup: func(repo MockcartGetter) {
				repo.EXPECT().
					Get(gomock.Any()).
					Return(&cart.Cart{
						Items: []cart.Item{
							{Sku: 1, Name: "1", Count: 1, Price: 100},
							{Sku: 3, Name: "3", Count: 1, Price: 100},
							{Sku: 2, Name: "2", Count: 1, Price: 100},
						},
						TotalPrice: 300,
					}, nil)
			},
			wantСart: cart.Cart{
				Items: []cart.Item{
					{Sku: 1, Name: "1", Count: 1, Price: 100},
					{Sku: 2, Name: "2", Count: 1, Price: 100},
					{Sku: 3, Name: "3", Count: 1, Price: 100},
				},
				TotalPrice: 300,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func(repo MockcartGetter) {
				repo.EXPECT().
					Get(gomock.Any()).
					Return(nil, cart.ErrProductNotFound)
			},
			wantErr: cart.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockcartGetter(ctrl)
			tt.setup(*mockRepo)
			uc := get.NewGetUseCase(mockRepo)

			c, err := uc.GetCart(userUd)
			if err != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.Equal(t, tt.wantСart, c)
			}
		})
	}
}
