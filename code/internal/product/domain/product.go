package domain

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Price       float64
	SalePrice   *float64
}

type UpdateProduct struct {
	Name        *string
	Description *string
	Price       *float64
	SalePrice   *float64
}

func NewProduct(name string, description *string, price float64, salePrice *float64) (Product, error) {
	if name == "" {
		return Product{}, NewValidationError(ErrValidationRequired, "name is required")
	}
	if price <= 0 {
		return Product{}, NewValidationError(ErrValidationRange, "price must be greater than 0")
	}
	if salePrice != nil && *salePrice < 0 {
		return Product{}, NewValidationError(ErrValidationRange, "sale_price must be greater than or equal to 0")
	}
	return Product{
		Name:        name,
		Description: description,
		Price:       price,
		SalePrice:   salePrice,
	}, nil
}

func NewUpdateProduct(name, description *string, price, salePrice *float64) (UpdateProduct, error) {
	if name != nil && *name == "" {
		return UpdateProduct{}, NewValidationError(ErrValidationRequired, "name must not be empty")
	}
	if price != nil && *price <= 0 {
		return UpdateProduct{}, NewValidationError(ErrValidationRange, "price must be greater than 0")
	}
	if salePrice != nil && *salePrice < 0 {
		return UpdateProduct{}, NewValidationError(ErrValidationRange, "sale_price must be greater than or equal to 0")
	}
	return UpdateProduct{
		Name:        name,
		Description: description,
		Price:       price,
		SalePrice:   salePrice,
	}, nil
}
