package usecase

import (
	portin "primo/service/internal/product/port/in"
	portout "primo/service/internal/product/port/out"
)

type Service struct {
	productRepo portout.ProductRepository
}

func NewService(productRepo portout.ProductRepository) portin.ProductService {
	return &Service{
		productRepo: productRepo,
	}
}
