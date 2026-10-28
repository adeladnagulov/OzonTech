package get

import "hw_2/internal/cart/domain/cart"

type cartGetter interface {
	Get(userId int64) (*cart.Cart, error)
}
