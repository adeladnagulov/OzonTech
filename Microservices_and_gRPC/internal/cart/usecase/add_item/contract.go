package additem

import "hw_2/internal/cart/domain/cart"

type cartCreator interface {
	Create(userId int64) *cart.Cart
}
