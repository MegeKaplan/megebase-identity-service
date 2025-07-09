package handlers

import (
	"log"
	"net/http"

	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.RegisterRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json",
				"error":   err.Error(),
			})
			return
		}

		user := models.User{
			Email:    body.Email,
			Password: body.Password,
		}

		log.Printf("email: %s, password: %s", user.Email, user.Password)

		c.JSON(http.StatusOK, dto.RegisterResponse{
			Token: "jwt",
			User:  user,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.LoginRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid json",
				"error":   err.Error(),
			})
			return
		}

		user := models.User{
			Email:    body.Email,
			Password: body.Password,
		}

		log.Printf("email: %s, password: %s", user.Email, user.Password)

		c.JSON(http.StatusOK, dto.LoginResponse{
			Token: "jwt",
			User:  user,
		})
	}
}
