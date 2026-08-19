package validation

import (
	"primo/service/internal/product/application"
	"primo/service/internal/product/handler/response"

	"github.com/google/uuid"
)

func ValidateUUID(id string) error {
	if id == "" {
		return response.NewValidationError(response.ErrValidationRequired, "id is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return response.NewValidationError(response.ErrValidationFormat, "id must be a valid UUID")
	}
	return nil
}

func ValidateCreateProductDto(dto application.CreateProductDto) error {
	if dto.Name == "" {
		return response.NewValidationError(response.ErrValidationRequired, "name is required")
	}

	if dto.Price <= 0 {
		return response.NewValidationError(response.ErrValidationRequired, "price must be greater than 0")
	}
	if dto.SalePrice != nil && *dto.SalePrice < 0 {
		return response.NewValidationError(response.ErrValidationRequired, "sale_price must be greater than or equal to 0")
	}
	return nil
}

func ValidatePatchProductDto(dto application.PatchProductDto) error {
	if dto.Name != nil && *dto.Name == "" {
		return response.NewValidationError(response.ErrValidationRequired, "name must not be empty")
	}
	if dto.Price != nil && *dto.Price <= 0 {
		return response.NewValidationError(response.ErrValidationRange, "price must be greater than 0")
	}
	if dto.SalePrice != nil && *dto.SalePrice < 0 {
		return response.NewValidationError(response.ErrValidationRange, "sale_price must be greater than or equal to 0")
	}
	return nil
}
