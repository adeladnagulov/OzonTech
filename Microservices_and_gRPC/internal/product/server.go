package product

import (
	"encoding/json"
	"fmt"
	pb "hw_2/pkg/api"
	"os"

	"github.com/go-playground/validator"
)

type ProductService struct {
	pb.UnimplementedProductServer

	products map[int64]*Product
	validate *validator.Validate
}

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func NewProductService(file string) (*ProductService, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read products file: %w", err)
	}
	var products map[int64]*Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, fmt.Errorf("failed to parse products JSON: %w", err)
	}

	return &ProductService{
		products: products,
		validate: validator.New(),
	}, nil
}
