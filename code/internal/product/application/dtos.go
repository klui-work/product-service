package application

type CreateProductDto struct {
	Name        string   `json:"name" example:"AirPods Pro 2"`
	Description *string  `json:"description" example:"Apple AirPods Pro 2nd Generation"`
	Price       float64  `json:"price" example:"8990"`
	SalePrice   *float64 `json:"sale_price" example:"7990"`
}

type PatchProductDto struct {
	Name        *string  `json:"name" example:"iPhone 16 Pro Max"`
	Description *string  `json:"description" example:"Apple iPhone 16 Pro Max 512GB"`
	Price       *float64 `json:"price" example:"52900"`
	SalePrice   *float64 `json:"sale_price" example:"49900"`
}

type ProductDto struct {
	ID          string   `json:"id" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567801"`
	Name        string   `json:"name" example:"iPhone 16 Pro"`
	Description *string  `json:"description" example:"Apple iPhone 16 Pro 256GB"`
	Price       float64  `json:"price" example:"48900"`
	SalePrice   *float64 `json:"sale_price" example:"45900"`
}
