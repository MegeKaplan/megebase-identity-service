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

	ErrEmailAlreadyExists = &AppError{
		Code:       "EMAIL_EXISTS",
		Message:    "Email already exists",
		HTTPStatus: 400,
	}

	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		HTTPStatus: 401,
	}

	ErrEmailNotFound = &AppError{
		Code:       "EMAIL_NOT_FOUND",
		Message:    "Email not found",
		HTTPStatus: 404,
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
		HTTPStatus: 401,
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

	ErrDBOperation = &AppError{
		Code:       "DB_OPERATION_ERROR",
		Message:    "Database operation failed",
		HTTPStatus: 500,
	}
)

// Hashing Errors
var (
	ErrPasswordHashingFailed = &AppError{
		Code:       "PASSWORD_HASHING_ERROR",
		Message:    "Failed to hash password",
		HTTPStatus: 500,
	}
)

// Token Generation Errors
var (
	ErrTokenGenerationFailed = &AppError{
		Code:       "TOKEN_GENERATION_FAILED",
		Message:    "Failed to generate token",
		HTTPStatus: 500,
	}

	ErrInvalidSigningMethod = &AppError{
		Code:       "INVALID_SIGNING_METHOD",
		Message:    "Unexpected signing method",
		HTTPStatus: 401,
	}

	ErrInvalidToken = &AppError{
		Code:       "INVALID_TOKEN",
		Message:    "Token is invalid or expired",
		HTTPStatus: 401,
	}
)
