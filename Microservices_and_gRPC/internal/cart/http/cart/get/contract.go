package get

import (
	"context"
	"hw_2/internal/cart/domain/cart"
)

type CartUseCase interface {
	GetCart(ctx context.Context, userId int64) (cart.Cart, error)
}
