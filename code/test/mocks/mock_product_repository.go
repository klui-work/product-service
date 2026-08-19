package mocks

import (
	"context"
	"primo/service/internal/product/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) PatchProduct(ctx context.Context, id uuid.UUID, update domain.UpdateProduct) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}
