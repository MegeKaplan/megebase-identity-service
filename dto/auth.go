package dto

import "github.com/MegeKaplan/megebase-identity-service/models"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"min=6"`
	OTP      string `json:"otp" binding:"required"`
}

type RegisterResponse struct {
	AccessToken string      `json:"access_token"`
	User        models.User `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        models.User `json:"user"`
}

type RefreshTokensResponse struct {
	AccessToken string `json:"access_token"`
}
