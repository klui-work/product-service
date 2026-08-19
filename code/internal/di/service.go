package di

import (
	portin "primo/service/internal/product/port/in"
	portout "primo/service/internal/product/port/out"
	"primo/service/internal/product/usecase"
)

func ProvideProductService(
	productRepo portout.ProductRepository,
) portin.ProductService {
	return usecase.NewService(productRepo)
}
