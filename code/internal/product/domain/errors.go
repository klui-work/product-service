package domain

import "fmt"

const (
	ErrValidationRequired = "E_VALIDATION_REQUIRED"
	ErrValidationFormat   = "E_VALIDATION_FORMAT"
	ErrValidationRange    = "E_VALIDATION_RANGE"
	ErrBadRequest         = "E_BAD_REQUEST"

	ErrUnauthorized = "E_UNAUTHORIZED"
	ErrForbidden    = "E_FORBIDDEN"

	ErrNotFound  = "E_NOT_FOUND"
	ErrDuplicate = "E_DUPLICATE"
	ErrConflict  = "E_CONFLICT"

	ErrDBConnection = "E_DB_CONNECTION"
	ErrDBQuery      = "E_DB_QUERY"
	ErrDBTimeout    = "E_DB_TIMEOUT"

	ErrInternal       = "E_INTERNAL"
	ErrServiceUnavail = "E_SERVICE_UNAVAILABLE"
)

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewValidationError(code string, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func NewDBError(code string, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func NewAuthError(code string, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func NewInternalError(msg string) *AppError {
	return &AppError{Code: ErrInternal, Message: msg}
}
