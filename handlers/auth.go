package handlers

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/services"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.RegisterRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		isValid, err := h.authService.VerifyOTP(body.Email, body.OTP)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		if !isValid {
			utils.JSONError(c, response.ErrInvalidOTP, "")
			return
		}

		user, err := h.authService.RegisterUser(body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		accessToken, refreshToken, err := h.authService.GenerateTokens(user)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		c.SetCookie(
			"refresh_token",
			refreshToken,
			7*24*60*60,
			"/",
			"",
			true,
			true,
		)

		utils.JSONSuccess(c, response.UserRegistered, dto.RegisterResponse{
			AccessToken: accessToken,
			User:        user,
		})
	}
}

func (h *authHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.LoginRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		user, err := h.authService.LoginUser(body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		accessToken, refreshToken, err := h.authService.GenerateTokens(user)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		c.SetCookie(
			"refresh_token",
			refreshToken,
			7*24*60*60,
			"/",
			"",
			true,
			true,
		)

		utils.JSONSuccess(c, response.UserLoggedIn, dto.LoginResponse{
			AccessToken: accessToken,
			User:        user,
		})
	}
}

func (h *authHandler) SendOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.SendOTPRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		if err := h.authService.SendOTP(body); err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.OTPSent, nil)
	}
}

func (h *authHandler) VerifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.VerifyOTPRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		isValid, err := h.authService.VerifyOTP(body.Email, body.OTP)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		if !isValid {
			utils.JSONError(c, response.ErrInvalidOTP, "")
			return
		}

		utils.JSONSuccess(c, response.OTPVerified, nil)
	}
}

func (h *authHandler) RefreshTokens() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			utils.JSONError(c, response.ErrMissingRefreshToken, "")
			return
		}

		accessToken, newRefreshToken, err := h.authService.RefreshTokens(refreshToken)
		if err.(*response.AppError) != nil {
			utils.JSONError(c, err.(*response.AppError), err.(*response.AppError).Details)
			return
		}

		c.SetCookie(
			"refresh_token",
			newRefreshToken,
			7*24*60*60,
			"/",
			"",
			true,
			true,
		)

		utils.JSONSuccess(c, response.TokensRefreshed, dto.RefreshTokensResponse{
			AccessToken: accessToken,
		})
	}
}
