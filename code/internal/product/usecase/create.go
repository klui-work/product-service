package usecase

import (
	"context"
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"
	"primo/service/internal/utils"
)

func (s *Service) CreateProduct(ctx context.Context, dto application.CreateProductDto) error {
	product, err := domain.NewProduct(dto.Name, dto.Description, dto.Price, dto.SalePrice)
	if err != nil {
		return err
	}
	product.ID = utils.GenerateUUID()
	return s.productRepo.CreateProduct(ctx, product)
}
