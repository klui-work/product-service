package di

import (
	v1 "primo/service/internal/product/handler/v1"

	"gorm.io/gorm"
)

func InitializeProductAdapter(db *gorm.DB) v1.IProductHandler {
	repository := ProvideProductRepository(db)
	service := ProvideProductService(repository)
	handler := ProvideProductHandler(service)
	return handler
}
