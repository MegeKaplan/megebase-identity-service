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
		HTTPStatus: http.StatusCreated,
		Code:       "USER_REGISTERED",
		Message:    "User registered successfully",
	}

	UserLoggedIn = &AppSuccess{
		HTTPStatus: http.StatusOK,
		Code:       "USER_LOGGED_IN",
		Message:    "User logged in successfully",
	}
)
