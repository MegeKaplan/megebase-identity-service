package response

import "net/http"

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
		HTTPStatus: http.StatusBadRequest,
	}

	ErrEmailAlreadyExists = &AppError{
		Code:       "EMAIL_EXISTS",
		Message:    "Email already exists",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidCredentials = &AppError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrEmailNotFound = &AppError{
		Code:       "EMAIL_NOT_FOUND",
		Message:    "Email not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}
)

// Database Errors
var (
	ErrDBConnection = &AppError{
		Code:       "DB_CONNECTION_ERROR",
		Message:    "Failed to connect to the database",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrDBMigration = &AppError{
		Code:       "DB_MIGRATION_ERROR",
		Message:    "Failed to migrate the database",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrDBOperation = &AppError{
		Code:       "DB_OPERATION_ERROR",
		Message:    "Database operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// Hashing Errors
var (
	ErrPasswordHashingFailed = &AppError{
		Code:       "PASSWORD_HASHING_ERROR",
		Message:    "Failed to hash password",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// Token Generation Errors
var (
	ErrTokenGenerationFailed = &AppError{
		Code:       "TOKEN_GENERATION_FAILED",
		Message:    "Failed to generate token",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrInvalidSigningMethod = &AppError{
		Code:       "INVALID_SIGNING_METHOD",
		Message:    "Unexpected signing method",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidToken = &AppError{
		Code:       "INVALID_TOKEN",
		Message:    "Token is invalid or expired",
		HTTPStatus: http.StatusUnauthorized,
	}
)
