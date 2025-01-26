package constants

const (
	// Database errors
	ErrDBConnection   = "failed to connect to database"
	ErrDBQuery        = "database query failed"
	ErrRecordNotFound = "record not found"

	// Validation errors
	ErrInvalidServiceID = "invalid service ID: must be a positive integer"
	ErrRequiredField    = "required field missing: %s"
	ErrInvalidFormat    = "invalid format for field: %s"

	// HTTP errors
	ErrInternalServer = "internal server error"
	ErrBadRequest     = "bad request"
	ErrUnauthorized   = "unauthorized access"
	ErrForbidden      = "forbidden access"
	ErrNotFound       = "resource not found"

	// Service specific errors
	ErrServiceNotFound     = "service not found"
	ErrServiceFetchFailed  = "Failed to fetch service"
	ErrVersionFetchFailed  = "Failed to fetch versions"
	ErrServiceCountFailed  = "Failed to count services"
	ErrServicesFetchFailed = "Failed to fetch services"
	ErrServiceDeleteFailed = "failed to delete service"
	ErrVersionNotFound     = "version not found"
)

type ServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
