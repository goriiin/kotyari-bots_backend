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
	ValidationMsg         = "VALIDATION_ERROR"
	NotFoundMsg           = "NOT_FOUND"
	ConflictMsg           = "CONFLICT"
	InternalMsg           = "INTERNAL_ERROR"
	NotImplementedMsg     = "NOT_IMPLEMENTED"
	ServiceUnavailableMsg = "SERVICE_UNAVAILABLE"
	UnauthorizedMsg       = "UNAUTHORIZED"
	InvalidMsg            = "INVALID_ERROR"
	RequiredMsg           = "REQUIRED_ERROR"
	MarshalMsg            = "MARSHALLING_ERROR"
	UnmarshalMsg          = "UNMARSHALLING_ERROR"
)
