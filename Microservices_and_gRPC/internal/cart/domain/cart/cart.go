package cart

type Cart struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

type Item struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type AddItemReq struct {
	Count uint32 `json:"count" validate:"required,gt=0"`
}
