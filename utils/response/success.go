package response

import "net/http"

type AppSuccess struct {
	HTTPStatus int         `json:"-"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// Auth
var (
	UserRegistered = &AppSuccess{
		Code:       "USER_REGISTERED",
		Message:    "User registered successfully",
		HTTPStatus: http.StatusCreated,
	}

	UserLoggedIn = &AppSuccess{
		Code:       "USER_LOGGED_IN",
		Message:    "User logged in successfully",
		HTTPStatus: http.StatusOK,
	}
)

// OTP
var (
	OTPSent = &AppSuccess{
		Code:       "OTP_SENT",
		Message:    "OTP sent successfully",
		HTTPStatus: http.StatusOK,
	}

	OTPVerified = &AppSuccess{
		Code:       "OTP_VERIFIED",
		Message:    "OTP verified successfully",
		HTTPStatus: http.StatusOK,
	}
)

// User
var (
	UserFetched = &AppSuccess{
		Code:       "USER_FETCHED",
		Message:    "User fetched successfully",
		HTTPStatus: http.StatusOK,
	}

	UsersFetched = &AppSuccess{
		Code:       "USERS_FETCHED",
		Message:    "Users fetched successfully",
		HTTPStatus: http.StatusOK,
	}

	UserUpdated = &AppSuccess{
		Code:       "USER_UPDATED",
		Message:    "User updated successfully",
		HTTPStatus: http.StatusOK,
	}

	UserDeleted = &AppSuccess{
		Code:       "USER_DELETED",
		Message:    "User deleted successfully",
		HTTPStatus: http.StatusOK,
	}
)

// Refresh Token
var (
	TokensRefreshed = &AppSuccess{
		Code:       "TOKENS_REFRESHED",
		Message:    "Tokens refreshed successfully",
		HTTPStatus: http.StatusOK,
	}
)
