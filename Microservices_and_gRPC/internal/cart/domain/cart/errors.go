package cart

import "errors"

var (
	ErrUserIdBadRequest = errors.New("user_id must be a natural number")
	ErrSkuIdBadRequest  = errors.New("sku_id must be a natural number")
	ErrBadRequest       = errors.New("bad request")
	ErrProductNotFound  = errors.New("product not found")
	ErrGRPCUnavailable  = errors.New("product-service unavailable")
)
