package repository

import (
	"hw_2/internal/cart/domain/cart"
	"sync"
)

type CartsRepo struct {
	data map[int64]*cart.Cart
	mu   sync.RWMutex
}

func NewCartsRepo() *CartsRepo {
	return &CartsRepo{
		data: make(map[int64]*cart.Cart),
	}
}

func (repo *CartsRepo) Update(userId int64, fn func(*cart.Cart)) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	c, exists := repo.data[userId]
	if !exists {
		c = &cart.Cart{
			Items:      []cart.Item{},
			TotalPrice: 0}
		repo.data[userId] = c
	}
	fn(c)
}

func (repo *CartsRepo) Create(userId int64) *cart.Cart {
	var rez *cart.Cart
	repo.Update(userId, func(c *cart.Cart) { rez = c })
	return rez
}

func (repo *CartsRepo) Get(userId int64) (*cart.Cart, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	c, exists := repo.data[userId]
	if !exists {
		return nil, cart.ErrProductNotFound
	}
	return c, nil
}
