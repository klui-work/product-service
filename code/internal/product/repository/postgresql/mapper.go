package postgresql

import (
	"primo/service/internal/product/domain"
)

func MapProductToEntity(p domain.Product) ProductEntity {
	return ProductEntity{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		SalePrice:   p.SalePrice,
	}
}

func MapEntityToProduct(e ProductEntity) domain.Product {
	return domain.Product{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Price:       e.Price,
		SalePrice:   e.SalePrice,
	}
}

func MapEntitiesToProducts(items []ProductEntity) []domain.Product {
	if len(items) == 0 {
		return []domain.Product{}
	}
	result := make([]domain.Product, len(items))
	for i, item := range items {
		result[i] = MapEntityToProduct(item)
	}
	return result
}

func MapUpdateProductToEntity(u domain.UpdateProduct) UpdateProductEntity {
	return UpdateProductEntity{
		Name:        u.Name,
		Description: u.Description,
		Price:       u.Price,
		SalePrice:   u.SalePrice,
	}
}
