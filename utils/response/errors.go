package response

type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Auth Errors
var (
	ErrInvalidJSON = &AppError{
		Code:       "INVALID_JSON",
		Message:    "Request body is not valid JSON",
		HTTPStatus: 400,
	}
)

// Database Errors
var (
	ErrDBConnection = &AppError{
		Code:       "DB_CONNECTION_ERROR",
		Message:    "Failed to connect to the database",
		HTTPStatus: 500,
	}

	ErrDBMigration = &AppError{
		Code:       "DB_MIGRATION_ERROR",
		Message:    "Failed to migrate the database",
		HTTPStatus: 500,
	}
)
