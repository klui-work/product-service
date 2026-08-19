package usecase

import (
	"context"
	"errors"
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"
	portin "primo/service/internal/product/port/in"
	productusecase "primo/service/internal/product/usecase"
	"primo/service/test/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(v string) *string    { return &v }
func f64Ptr(v float64) *float64 { return &v }

func setupService() (portin.ProductService, *mocks.MockProductRepository) {
	mockRepo := new(mocks.MockProductRepository)
	svc := productusecase.NewService(mockRepo)
	return svc, mockRepo
}

func TestCreateProduct_Success(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()

	mockRepo.On("CreateProduct", ctx, mock.AnythingOfType("domain.Product")).Return(nil)

	dto := application.CreateProductDto{
		Name:        "iPhone",
		Description: strPtr("Apple iPhone"),
		Price:       48900,
		SalePrice:   f64Ptr(45900),
	}

	err := svc.CreateProduct(ctx, dto)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "CreateProduct", ctx, mock.AnythingOfType("domain.Product"))
}

func TestCreateProduct_NullableFields(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()

	mockRepo.On("CreateProduct", ctx, mock.MatchedBy(func(p domain.Product) bool {
		return p.Name == "iPhone" && p.Description == nil && p.SalePrice == nil && p.Price == 48900
	})).Return(nil)

	dto := application.CreateProductDto{
		Name:  "iPhone",
		Price: 48900,
	}

	err := svc.CreateProduct(ctx, dto)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_RepoError(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()

	mockRepo.On("CreateProduct", ctx, mock.AnythingOfType("domain.Product")).
		Return(errors.New("db connection error"))

	dto := application.CreateProductDto{
		Name:  "iPhone",
		Price: 48900,
	}

	err := svc.CreateProduct(ctx, dto)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db connection error")
}

func TestCreateProduct_InvalidDomainRules(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()

	dto := application.CreateProductDto{
		Name:  "",
		Price: 48900,
	}

	err := svc.CreateProduct(ctx, dto)

	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.ErrValidationRequired, appErr.Code)
}

func TestGetAllProducts_Success(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()

	expected := []domain.Product{
		{ID: uuid.New(), Name: "iPhone", Price: 100},
	}
	mockRepo.On("GetAllProducts", ctx).Return(expected, nil)

	products, err := svc.GetAllProducts(ctx)

	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "iPhone", products[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct_Success(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("PatchProduct", ctx, id, mock.AnythingOfType("domain.UpdateProduct")).Return(nil)

	name := "Updated iPhone"
	dto := application.PatchProductDto{
		Name: &name,
	}

	err := svc.PatchProduct(ctx, id, dto)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "PatchProduct", ctx, id, mock.AnythingOfType("domain.UpdateProduct"))
}

func TestPatchProduct_RepoNotFound(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()
	id := uuid.New()

	mockRepo.On("PatchProduct", ctx, id, mock.AnythingOfType("domain.UpdateProduct")).
		Return(domain.NewValidationError(domain.ErrNotFound, "product not found"))

	name := "Updated"
	dto := application.PatchProductDto{
		Name: &name,
	}

	err := svc.PatchProduct(ctx, id, dto)

	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.ErrNotFound, appErr.Code)
}

func TestPatchProduct_PartialUpdate(t *testing.T) {
	svc, mockRepo := setupService()
	ctx := context.Background()
	id := uuid.New()

	price := 99.99
	dto := application.PatchProductDto{
		Price: &price,
	}

	mockRepo.On("PatchProduct", ctx, id, mock.MatchedBy(func(u domain.UpdateProduct) bool {
		return u.Price != nil && *u.Price == 99.99 && u.Name == nil && u.Description == nil && u.SalePrice == nil
	})).Return(nil)

	err := svc.PatchProduct(ctx, id, dto)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatchProduct_InvalidDomainRules(t *testing.T) {
	svc, _ := setupService()
	ctx := context.Background()
	id := uuid.New()

	name := ""
	dto := application.PatchProductDto{
		Name: &name,
	}

	err := svc.PatchProduct(ctx, id, dto)

	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.ErrValidationRequired, appErr.Code)
}
