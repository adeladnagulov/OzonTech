package cart

import (
	"context"
	"fmt"
	pb "hw_2/pkg/api"
	"sort"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	productClient pb.ProductClient
	carts         map[int64]*Cart
	validate      *validator.Validate
}

func NewCartService(productClient pb.ProductClient) *Service {
	return &Service{
		productClient: productClient,
		carts:         make(map[int64]*Cart),
		validate:      validator.New(),
	}
}

func (s *Service) AddItem(req AddItemReq, ctx context.Context, skuId, userId int64) error {
	err := s.validate.Struct(req)
	if err != nil {
		return fmt.Errorf("%w body (count): %v", ErrBadRequest, err)
	}

	product, err := s.productClient.GetProduct(ctx, &pb.GetProductRequest{Sku: skuId})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("%w: sku %d", ErrProductNotFound, skuId)
		}
		return fmt.Errorf("%w: %v", ErrGRPCUnavailable, err)
	}

	cart, exists := s.carts[userId]
	if !exists {
		s.carts[int64(userId)] = &Cart{
			Items:      []Item{},
			TotalPrice: 0,
		}
		cart = s.carts[int64(userId)]
	}

	if idx := findItem(cart, product.GetSku()); idx != -1 {
		cart.Items[idx].Count += req.Count
		cart.TotalPrice += req.Count * cart.Items[idx].Price
	} else {
		cart.Items = append(cart.Items, Item{
			Sku:   product.Sku,
			Name:  product.Name,
			Count: req.Count,
			Price: uint32(product.Price),
		})
		cart.TotalPrice += uint32(product.Price) * req.Count
	}
	return nil
}

func (s *Service) GetCart(userId int64) (Cart, error) {
	cart, exists := s.carts[userId]
	if !exists {
		return Cart{}, fmt.Errorf("%w: for user: %d", ErrProductNotFound, userId)
	}
	sortedItems := make([]Item, len(cart.Items))
	copy(sortedItems, cart.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return cart.Items[i].Sku < cart.Items[j].Sku
	})

	return Cart{
		Items:      sortedItems,
		TotalPrice: cart.TotalPrice,
	}, nil
}

func findItem(cart *Cart, sku int64) int {
	for idx, item := range cart.Items {
		if sku == item.Sku {
			return idx
		}
	}
	return -1
}
