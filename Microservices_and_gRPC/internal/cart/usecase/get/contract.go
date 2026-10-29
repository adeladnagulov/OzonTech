package get

import "hw_2/internal/cart/domain/cart"

//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test

type cartGetter interface {
	Get(userId int64) (*cart.Cart, error)
}
