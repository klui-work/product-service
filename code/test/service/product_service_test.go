package service

import (
	"primo/service/internal/product/application"
	"primo/service/internal/product/domain"
	"primo/service/internal/product/handler/mapper"
	"primo/service/internal/product/handler/validation"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func strPtr(v string) *string    { return &v }
func f64Ptr(v float64) *float64 { return &v }

func TestMapProductToDto(t *testing.T) {
	id := uuid.New()
	product := domain.Product{
		ID:          id,
		Name:        "iPhone",
		Description: strPtr("Apple iPhone"),
		Price:       48900,
		SalePrice:   f64Ptr(45900),
	}

	resp := mapper.MapProductToDto(product)

	assert.Equal(t, id.String(), resp.ID)
	assert.Equal(t, "iPhone", resp.Name)
	assert.Equal(t, 48900.0, resp.Price)
	assert.Equal(t, f64Ptr(45900), resp.SalePrice)
}

func TestNewProduct_Success(t *testing.T) {
	product, err := domain.NewProduct("iPhone", strPtr("Apple iPhone"), 48900, f64Ptr(45900))

	assert.NoError(t, err)
	assert.Equal(t, "iPhone", product.Name)
	assert.Equal(t, 48900.0, product.Price)
}

func TestNewProduct_NullableFields(t *testing.T) {
	product, err := domain.NewProduct("iPhone", nil, 48900, nil)

	assert.NoError(t, err)
	assert.Nil(t, product.Description)
	assert.Nil(t, product.SalePrice)
}

func TestNewProduct_EmptyName(t *testing.T) {
	_, err := domain.NewProduct("", strPtr("desc"), 100, f64Ptr(80))

	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.ErrValidationRequired, appErr.Code)
}

func TestNewProduct_ZeroPrice(t *testing.T) {
	_, err := domain.NewProduct("iPhone", strPtr("desc"), 0, f64Ptr(80))

	assert.Error(t, err)
	appErr, ok := err.(*domain.AppError)
	assert.True(t, ok)
	assert.Equal(t, domain.ErrValidationRange, appErr.Code)
}

func TestNewUpdateProduct_PartialFields(t *testing.T) {
	name := "Only Name"

	update, err := domain.NewUpdateProduct(&name, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, &name, update.Name)
	assert.Nil(t, update.Description)
}

func TestNewUpdateProduct_EmptyName(t *testing.T) {
	name := ""

	_, err := domain.NewUpdateProduct(&name, nil, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name must not be empty")
}

func TestValidateCreateProductDto_Success(t *testing.T) {
	dto := application.CreateProductDto{
		Name:      "iPhone",
		Price:     48900,
		SalePrice: f64Ptr(45900),
	}

	err := validation.ValidateCreateProductDto(dto)
	assert.NoError(t, err)
}

func TestValidateCreateProductDto_OmitNullableFields(t *testing.T) {
	dto := application.CreateProductDto{
		Name:  "iPhone",
		Price: 48900,
	}

	err := validation.ValidateCreateProductDto(dto)
	assert.NoError(t, err)
}

func TestValidateCreateProductDto_EmptyName(t *testing.T) {
	dto := application.CreateProductDto{
		Name:  "",
		Price: 48900,
	}

	err := validation.ValidateCreateProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestValidateCreateProductDto_ZeroPrice(t *testing.T) {
	dto := application.CreateProductDto{
		Name:  "iPhone",
		Price: 0,
	}

	err := validation.ValidateCreateProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price must be greater than 0")
}

func TestValidateCreateProductDto_NegativeSalePrice(t *testing.T) {
	dto := application.CreateProductDto{
		Name:      "iPhone",
		Price:     100,
		SalePrice: f64Ptr(-1),
	}

	err := validation.ValidateCreateProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sale_price")
}

func TestValidatePatchProductDto_Success(t *testing.T) {
	name := "Updated"
	price := 100.0

	dto := application.PatchProductDto{
		Name:  &name,
		Price: &price,
	}

	err := validation.ValidatePatchProductDto(dto)
	assert.NoError(t, err)
}

func TestValidatePatchProductDto_EmptyName(t *testing.T) {
	name := ""

	dto := application.PatchProductDto{
		Name: &name,
	}

	err := validation.ValidatePatchProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name must not be empty")
}

func TestValidatePatchProductDto_NegativePrice(t *testing.T) {
	price := -10.0

	dto := application.PatchProductDto{
		Price: &price,
	}

	err := validation.ValidatePatchProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price must be greater than 0")
}

func TestValidatePatchProductDto_NegativeSalePrice(t *testing.T) {
	salePrice := -5.0

	dto := application.PatchProductDto{
		SalePrice: &salePrice,
	}

	err := validation.ValidatePatchProductDto(dto)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sale_price")
}

func TestValidatePatchProductDto_EmptyBody(t *testing.T) {
	dto := application.PatchProductDto{}

	err := validation.ValidatePatchProductDto(dto)
	assert.NoError(t, err)
}
