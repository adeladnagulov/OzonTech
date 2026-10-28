package repository

import "hw_2/internal/cart/domain/cart"

type CartsRepo struct {
	data map[int64]*cart.Cart
}

func NewCartsRepo() *CartsRepo {
	return &CartsRepo{
		data: make(map[int64]*cart.Cart),
	}
}

func (repo *CartsRepo) Create(userId int64) *cart.Cart {
	c, exists := repo.data[userId]
	if !exists {
		repo.data[int64(userId)] = &cart.Cart{
			Items:      []cart.Item{},
			TotalPrice: 0,
		}
		c = repo.data[int64(userId)]
	}
	return c
}

func (repo *CartsRepo) Get(userId int64) (*cart.Cart, error) {
	c, exists := repo.data[userId]
	if !exists {
		return nil, cart.ErrProductNotFound
	}
	return c, nil
}
