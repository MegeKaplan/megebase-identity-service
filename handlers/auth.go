package handlers

import (
	"log"

	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/models"
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

		user := models.User{
			Email:    body.Email,
			Password: body.Password,
		}

		log.Printf("email: %s, password: %s", user.Email, user.Password)

		utils.JSONSuccess(c, response.UserRegistered, dto.RegisterResponse{
			Token: "jwt",
			User:  user,
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

		user := models.User{
			Email:    body.Email,
			Password: body.Password,
		}

		log.Printf("email: %s, password: %s", user.Email, user.Password)

		utils.JSONSuccess(c, response.UserLoggedIn, dto.LoginResponse{
			Token: "jwt",
			User:  user,
		})
	}
}
