package additem

import (
	"context"
	"fmt"

	pb "hw_2/pkg/api"

	"hw_2/internal/cart/domain/cart"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddUseCase struct {
	repo          cartCreator
	productClient cartProductClient
	validate      *validator.Validate
}

func NewAddUseCase(repo cartCreator, productClient cartProductClient) *AddUseCase {
	return &AddUseCase{
		repo:          repo,
		productClient: productClient,
		validate:      validator.New(),
	}
}

func (uc *AddUseCase) AddItem(req cart.AddItemReq, ctx context.Context, skuId, userId int64) error {
	err := uc.validate.Struct(req)
	if err != nil {
		return fmt.Errorf("%w body (count): %v", cart.ErrBadRequest, err)
	}

	product, err := uc.productClient.GetProduct(ctx, &pb.GetProductRequest{Sku: skuId})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("%w: sku %d", cart.ErrProductNotFound, skuId)
		}
		return fmt.Errorf("%w: %v", cart.ErrGRPCUnavailable, err)
	}

	c := uc.repo.Create(userId)

	if idx := findItem(c, product.GetSku()); idx != -1 {
		c.Items[idx].Count += req.Count
		c.TotalPrice += req.Count * c.Items[idx].Price
	} else {
		c.Items = append(c.Items, cart.Item{
			Sku:   product.Sku,
			Name:  product.Name,
			Count: req.Count,
			Price: uint32(product.Price),
		})
		c.TotalPrice += uint32(product.Price) * req.Count
	}
	return nil
}

func findItem(cart *cart.Cart, sku int64) int {
	for idx, item := range cart.Items {
		if sku == item.Sku {
			return idx
		}
	}
	return -1
}
