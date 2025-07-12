package handlers

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/services"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.RegisterRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		createdUser, err := services.RegisterUser(db, body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.UserRegistered, dto.RegisterResponse{
			Token: "jwt",
			User:  createdUser,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.LoginRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			utils.JSONError(c, response.ErrInvalidJSON, err.Error())
			return
		}

		existingUser, err := services.LoginUser(db, body)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.UserLoggedIn, dto.LoginResponse{
			Token: "jwt",
			User:  existingUser,
		})
	}
}
