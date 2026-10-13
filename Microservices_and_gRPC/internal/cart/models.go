package cart

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

type AddItemReq struct {
	Count uint32 `json:"count" validate:"required,gt=0"`
}
