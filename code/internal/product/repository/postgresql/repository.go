package postgresql

import (
	"context"
	"primo/service/internal/product/domain"
	"primo/service/internal/product/port/out"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var table = "tb_products"

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) out.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var productEntities []ProductEntity
	if err := r.db.WithContext(ctx).Table(table).Find(&productEntities).Error; err != nil {
		return nil, domain.NewDBError(domain.ErrDBQuery, err.Error())
	}
	return MapEntitiesToProducts(productEntities), nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entity := MapProductToEntity(product)
	if err := r.db.WithContext(ctx).Table(table).Create(&entity).Error; err != nil {
		return domain.NewDBError(domain.ErrDBQuery, err.Error())
	}
	return nil
}

func (r *ProductRepository) PatchProduct(ctx context.Context, id uuid.UUID, update domain.UpdateProduct) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entity := MapUpdateProductToEntity(update)
	result := r.db.WithContext(ctx).Table(table).Where("id = ?", id).Updates(entity)
	if result.Error != nil {
		return domain.NewDBError(domain.ErrDBQuery, result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return domain.NewValidationError(domain.ErrNotFound, "product not found")
	}
	return nil
}
