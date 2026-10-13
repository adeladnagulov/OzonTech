package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	pb "hw_2/pkg/api"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService struct {
	productClient pb.ProductClient
	carts         map[int64]*Cart
	validate      *validator.Validate
}

func NewCartService(productClient pb.ProductClient) *CartService {
	return &CartService{
		productClient: productClient,
		carts:         make(map[int64]*Cart),
		validate:      validator.New(),
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

var (
	ErrUserIdBadRequest = errors.New("user_id must be a natural number")
	ErrSkuIdBadRequest  = errors.New("sku_id must be a natural number")
)

type AddItemReq struct {
	Count uint32 `json:"count" validate:"required,gt=0"`
}

func (s *CartService) AddItem(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed parse user_id: %w", err))
		return
	}
	if userId <= 0 {
		writeError(w, http.StatusBadRequest, ErrUserIdBadRequest)
		return
	}
	skuIdStr := r.PathValue("sku_id")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed parse sku_id: %w", err))
		return
	}
	if skuId <= 0 {
		writeError(w, http.StatusBadRequest, ErrSkuIdBadRequest)
		return
	}

	var req AddItemReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed decode addItem body: %w", err))
		return
	}
	err = s.validate.Struct(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("failed validate req body: %w", err))
		return
	}

	product, err := s.productClient.GetProduct(r.Context(), &pb.GetProductRequest{Sku: int64(skuId)})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			writeError(w, http.StatusPreconditionFailed, fmt.Errorf("product with such an ID does not exist: %w", err))
			return
		}
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed getting product: %w", err))
		return
	}

	cart, exists := s.carts[int64(userId)]
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
	w.WriteHeader(http.StatusOK)
}

func (s *CartService) GetCart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func writeError(w http.ResponseWriter, code int, desc error) {
	var problem struct {
		Code int    `json:"code"`
		Desc string `json:"description"`
	}
	problem.Code = code
	problem.Desc = desc.Error()
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(problem)
}

func findItem(cart *Cart, sku int64) int {
	for idx, item := range cart.Items {
		if sku == item.Sku {
			return idx
		}
	}
	return -1
}
