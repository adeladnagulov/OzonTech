package get

import (
	"fmt"
	"hw_2/internal/cart/domain/cart"
	pb "hw_2/pkg/api"
	"sort"
)

type GetUseCase struct {
	repo          cartGetter
	productClient pb.ProductClient
}

func NewGetUseCase(repo cartGetter, productClient pb.ProductClient) *GetUseCase {
	return &GetUseCase{
		repo:          repo,
		productClient: productClient,
	}
}

func (uc *GetUseCase) GetCart(userId int64) (cart.Cart, error) {
	c, err := uc.repo.Get(userId)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("%w: for user: %d", err, userId)
	}
	sortedItems := make([]cart.Item, len(c.Items))
	copy(sortedItems, c.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return c.Items[i].Sku < c.Items[j].Sku
	})

	return cart.Cart{
		Items:      sortedItems,
		TotalPrice: c.TotalPrice,
	}, nil
}
