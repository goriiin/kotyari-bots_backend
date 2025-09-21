package constants

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

const (
	TaskStatusPending = "pending"
)

const (
	ErrValidationMsg         = "VALIDATION_ERROR"
	ErrNotFoundMsg           = "NOT_FOUND"
	ErrConflictMsg           = "CONFLICT"
	ErrInternalMsg           = "INTERNAL_ERROR"
	ErrNotImplementedMsg     = "NOT_IMPLEMENTED"
	ErrServiceUnavailableMsg = "SERVICE_UNAVAILABLE"
	ErrUnauthorizedMsg       = "UNAUTHORIZED"
	ErrInvalidMsg            = "INVALID_ERROR"
	ErrRequiredMsg           = "REQUIRED_ERROR"
)
