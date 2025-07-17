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

		createdUser, err := h.authService.RegisterUser(body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		token, err := utils.GenerateJWT(createdUser.ID, createdUser.Email)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.UserRegistered, dto.RegisterResponse{
			Token: token,
			User:  createdUser,
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

		existingUser, err := h.authService.LoginUser(body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		token, err := utils.GenerateJWT(existingUser.ID, existingUser.Email)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.UserLoggedIn, dto.LoginResponse{
			Token: token,
			User:  existingUser,
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
