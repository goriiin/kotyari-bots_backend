package constants

import "errors"

var (
	ErrNotFound           = errors.New(NotFoundMsg)
	ErrInvalid            = errors.New(InvalidMsg)
	ErrRequired           = errors.New(RequiredMsg)
	ErrValidation         = errors.New(ValidationMsg)
	ErrServiceUnavailable = errors.New(ServiceUnavailableMsg)
	ErrConflict           = errors.New(ConflictMsg)
	ErrInternal           = errors.New(InternalMsg)
	ErrUnauthorized       = errors.New(UnauthorizedMsg)
)
