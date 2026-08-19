package usecase

import (
	"context"
	"primo/service/internal/product/domain"
)

func (s *Service) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.productRepo.GetAllProducts(ctx)
}
