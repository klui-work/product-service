package in

import (
	"context"
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"

	"github.com/google/uuid"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	CreateProduct(ctx context.Context, dto application.CreateProductDto) error
	PatchProduct(ctx context.Context, id uuid.UUID, dto application.PatchProductDto) error
}
