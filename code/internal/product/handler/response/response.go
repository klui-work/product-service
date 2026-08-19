package response

import (
	"net/http"

	"primo/service/internal/product/domain"

	"github.com/gofiber/fiber/v2"
)

// Response represents the standard API response
// @Description Standard API response format
type Response struct {
	Successful bool   `json:"successful" example:"true"`
	ErrorCode  string `json:"error_code" example:""`
	Message    string `json:"message" example:""`
	Data       any    `json:"data" swaggertype:"object"`
}

// PatchResponse represents the PATCH /product/{id} success response
// @Description PATCH product success response format
type PatchResponse struct {
	Successful bool   `json:"successful" example:"true"`
	ErrorCode  string `json:"error_code" example:""`
}

// ErrorResponse represents the standard error response
// @Description Standard API error response format
type ErrorResponse struct {
	Successful bool   `json:"successful" example:"false"`
	ErrorCode  string `json:"error_code" example:"E_VALIDATION_REQUIRED"`
	Message    string `json:"message" example:"name is required"`
	Data       any    `json:"data" swaggertype:"object"`
}

var emptyData = map[string]any{}

var statusMap = map[string]int{
	domain.ErrValidationRequired: http.StatusBadRequest,
	domain.ErrValidationFormat:   http.StatusBadRequest,
	domain.ErrValidationRange:    http.StatusBadRequest,
	domain.ErrBadRequest:         http.StatusBadRequest,
	domain.ErrUnauthorized:       http.StatusUnauthorized,
	domain.ErrForbidden:          http.StatusForbidden,
	domain.ErrNotFound:           http.StatusNotFound,
	domain.ErrDuplicate:          http.StatusConflict,
	domain.ErrConflict:           http.StatusConflict,
	domain.ErrDBConnection:       http.StatusInternalServerError,
	domain.ErrDBQuery:            http.StatusInternalServerError,
	domain.ErrDBTimeout:          http.StatusGatewayTimeout,
	domain.ErrInternal:           http.StatusInternalServerError,
	domain.ErrServiceUnavail:     http.StatusServiceUnavailable,
}

func statusCodeFor(err *domain.AppError) int {
	if code, ok := statusMap[err.Code]; ok {
		return code
	}
	return http.StatusInternalServerError
}

func Success(c *fiber.Ctx, data ...any) error {
	d := any(emptyData)
	if len(data) > 0 && data[0] != nil {
		d = data[0]
	}
	return c.Status(http.StatusOK).JSON(Response{
		Successful: true,
		ErrorCode:  "",
		Message:    "",
		Data:       d,
	})
}

func SuccessPatch(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(PatchResponse{
		Successful: true,
		ErrorCode:  "",
	})
}

func Error(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*domain.AppError); ok {
		return c.Status(statusCodeFor(appErr)).JSON(Response{
			Successful: false,
			ErrorCode:  appErr.Code,
			Message:    appErr.Message,
			Data:       emptyData,
		})
	}
	return c.Status(http.StatusInternalServerError).JSON(Response{
		Successful: false,
		ErrorCode:  domain.ErrInternal,
		Message:    err.Error(),
		Data:       emptyData,
	})
}
