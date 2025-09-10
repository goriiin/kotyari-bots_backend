package constants

import "errors"

var (
	ErrNotFound   = errors.New(ErrNotFoundMsg)
	ErrInvalid    = errors.New(ErrInvalidMsg)
	ErrRequired   = errors.New(ErrRequiredMsg)
	ErrValidation = errors.New(ErrValidationMsg)
)
