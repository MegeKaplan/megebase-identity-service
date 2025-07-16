package handlers

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/services"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/gin-gonic/gin"
)

func Register(userRepo repositories.UserRepository, otpRepo repositories.OTPRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.RegisterRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		isValid, err := services.VerifyOTP(otpRepo, body.Email, body.OTP)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		if !isValid {
			utils.JSONError(c, response.ErrInvalidOTP, "")
			return
		}

		createdUser, err := services.RegisterUser(userRepo, body)
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

func Login(repo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.LoginRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		existingUser, err := services.LoginUser(repo, body)
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

func SendOTP(repo repositories.OTPRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.SendOTPRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		if err := services.SendOTP(repo, body); err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.OTPSent, nil)
	}
}

func VerifyOTP(repo repositories.OTPRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.VerifyOTPRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		isValid, err := services.VerifyOTP(repo, body.Email, body.OTP)
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
