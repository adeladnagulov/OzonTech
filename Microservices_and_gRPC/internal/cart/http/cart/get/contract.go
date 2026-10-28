package get

import "hw_2/internal/cart/domain/cart"

type CartUseCase interface {
	GetCart(userId int64) (cart.Cart, error)
}
