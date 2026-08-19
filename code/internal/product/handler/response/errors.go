package response

import "primo/service/internal/product/domain"

const (
	ErrValidationRequired = domain.ErrValidationRequired
	ErrValidationFormat   = domain.ErrValidationFormat
	ErrValidationRange    = domain.ErrValidationRange
	ErrBadRequest         = domain.ErrBadRequest
	ErrUnauthorized       = domain.ErrUnauthorized
	ErrForbidden          = domain.ErrForbidden
	ErrNotFound           = domain.ErrNotFound
	ErrDuplicate          = domain.ErrDuplicate
	ErrConflict           = domain.ErrConflict
	ErrDBConnection       = domain.ErrDBConnection
	ErrDBQuery            = domain.ErrDBQuery
	ErrDBTimeout          = domain.ErrDBTimeout
	ErrInternal           = domain.ErrInternal
	ErrServiceUnavail     = domain.ErrServiceUnavail
)

func NewValidationError(code string, msg string) *domain.AppError {
	return domain.NewValidationError(code, msg)
}

func NewDBError(code string, msg string) *domain.AppError {
	return domain.NewDBError(code, msg)
}

func NewAuthError(code string, msg string) *domain.AppError {
	return domain.NewAuthError(code, msg)
}

func NewInternalError(msg string) *domain.AppError {
	return domain.NewInternalError(msg)
}
