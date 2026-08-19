package out

import (
	"context"
	"primo/service/internal/product/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
	CreateProduct(ctx context.Context, product domain.Product) error
	PatchProduct(ctx context.Context, id uuid.UUID, update domain.UpdateProduct) error
}
