package additem

import (
	"context"
	"hw_2/internal/cart/domain/cart"
)

type CartUseCase interface {
	AddItem(req cart.AddItemReq, ctx context.Context, skuId, userId int64) error
}
