package response

import "primo/service/internal/product/application"

// ProductListResponse represents GET /product success response.
type ProductListResponse struct {
	Successful bool                     `json:"successful" example:"true"`
	ErrorCode  string                   `json:"error_code" example:""`
	Message    string                   `json:"message" example:""`
	Data       []application.ProductDto `json:"data"`
}

// HealthCheckData represents GET /product/health payload.
type HealthCheckData struct {
	Msg string `json:"msg" example:"Product API is healthy"`
}

// HealthCheckResponse represents GET /product/health success response.
type HealthCheckResponse struct {
	Successful bool            `json:"successful" example:"true"`
	ErrorCode  string          `json:"error_code" example:""`
	Message    string          `json:"message" example:""`
	Data       HealthCheckData `json:"data"`
}

// CreateProductResponse represents POST /product success response.
type CreateProductResponse struct {
	Successful bool                   `json:"successful" example:"true"`
	ErrorCode  string                 `json:"error_code" example:""`
	Message    string                 `json:"message" example:""`
	Data       map[string]interface{} `json:"data"`
}
