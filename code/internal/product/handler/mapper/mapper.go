package mapper

import (
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"
)

func MapProductToDto(p domain.Product) application.ProductDto {
	return application.ProductDto{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SalePrice:   p.SalePrice,
	}
}

func MapProductsToDtos(products []domain.Product) []application.ProductDto {
	if len(products) == 0 {
		return []application.ProductDto{}
	}
	result := make([]application.ProductDto, len(products))
	for i, p := range products {
		result[i] = MapProductToDto(p)
	}
	return result
}
