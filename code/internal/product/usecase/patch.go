package usecase

import (
	"context"
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"

	"github.com/google/uuid"
)

func (s *Service) PatchProduct(ctx context.Context, id uuid.UUID, dto application.PatchProductDto) error {
	update, err := domain.NewUpdateProduct(dto.Name, dto.Description, dto.Price, dto.SalePrice)
	if err != nil {
		return err
	}
	return s.productRepo.PatchProduct(ctx, id, update)
}
