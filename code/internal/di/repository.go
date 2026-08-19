package di

import (
	"primo/service/internal/product/repository/postgresql"
	portout "primo/service/internal/product/port/out"

	"gorm.io/gorm"
)

func ProvideProductRepository(
	db *gorm.DB,
) portout.ProductRepository {
	return postgresql.NewProductRepository(db)
}
