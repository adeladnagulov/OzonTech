package get

import (
	"context"
	"fmt"
	"hw_2/internal/cart/domain/cart"
	"sort"
	"time"

	pb "hw_2/pkg/api"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUseCase struct {
	repo          cartGetter
	productClient cartProductClient
}

func NewGetUseCase(repo cartGetter, productClient cartProductClient) *GetUseCase {
	return &GetUseCase{
		repo:          repo,
		productClient: productClient,
	}
}

func (uc *GetUseCase) GetCart(ctx context.Context, userId int64) (cart.Cart, error) {
	c, err := uc.repo.Get(userId)
	if err != nil {
		return cart.Cart{}, fmt.Errorf("%w: for user: %d", err, userId)
	}
	if len(c.Items) == 0 {
		return cart.Cart{Items: []cart.Item{}, TotalPrice: 0}, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	results := make([]*pb.GetProductResponse, len(c.Items))

	for i, item := range c.Items {
		g.Go(func() error {
			resp, err := uc.productClient.GetProduct(ctx, &pb.GetProductRequest{Sku: item.Sku})
			if err != nil {
				return fmt.Errorf("failed to get product %d: %w", item.Sku, err)
			}

			results[i] = resp
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		if status.Code(err) == codes.NotFound {
			return cart.Cart{}, cart.ErrProductNotFound
		}
		return cart.Cart{}, cart.ErrGRPCUnavailable
	}

	finalItems := make([]cart.Item, len(results))
	var totalPrice uint32

	for i, resp := range results {
		finalItems[i] = cart.Item{
			Sku:   resp.Sku,
			Name:  resp.Name,
			Count: c.Items[i].Count,
			Price: uint32(resp.Price),
		}
		totalPrice += uint32(resp.Price) * c.Items[i].Count
	}

	sort.Slice(finalItems, func(i, j int) bool {
		return finalItems[i].Sku < finalItems[j].Sku
	})

	return cart.Cart{
		Items:      finalItems,
		TotalPrice: totalPrice,
	}, nil
}
