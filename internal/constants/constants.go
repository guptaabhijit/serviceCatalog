package constants

const (
	// Database settings
	DefaultDBHost = "localhost"
	DefaultDBPort = 5432
	DefaultDBUser = "postgres"
	DefaultDBName = "servicecatalog"

	// Pagination settings
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultPage     = 1

	// Sort settings
	DefaultSortField = "id"
	DefaultSortOrder = "asc"
	DescSortOrder    = "desc"

	True        = "true"
	ShowDeleted = "showDeleted"
	Name        = "name"
	Description = "description"
	Error       = "error"
)
