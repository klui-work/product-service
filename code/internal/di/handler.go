package di

import (
	portin "primo/service/internal/product/port/in"
	"primo/service/internal/product/handler/v1"
)

func ProvideProductHandler(
	productService portin.ProductService,
) v1.IProductHandler {
	return v1.NewProductHandler(productService)
}
