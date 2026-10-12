package cart

import (
	pb "hw_2/pkg/api"
	"net/http"
)

type CartService struct {
	productClient pb.ProductClient
	carts         map[int64]*Cart
}

func NewCartService(productClient pb.ProductClient) *CartService {
	return &CartService{
		productClient: productClient,
		carts:         make(map[int64]*Cart),
	}
}

type Cart struct {
	Items      []Item
	TotalPrice uint32
}

type Item struct {
	Sku   int64
	Name  string
	Count uint32
	Price uint32
}

func (s *CartService) AddItem(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *CartService) GetCart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
