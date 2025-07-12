package services

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, body dto.RegisterRequest) (models.User, *response.AppError) {
	var existingUser models.User

	if err := db.First(&existingUser, "email = ?", body.Email).Error; err == nil {
		return models.User{}, response.ErrEmailAlreadyExists
	}

	user := models.User{
		Email:    body.Email,
		Password: body.Password,
	}

	if err := db.Create(&user).Error; err != nil {
		return models.User{}, response.ErrDBOperation
	}

	return user, nil
}

func LoginUser(db *gorm.DB, body dto.LoginRequest) (models.User, *response.AppError) {
	var existingUser models.User

	if err := db.First(&existingUser, "email = ?", body.Email).Error; err != nil {
		return models.User{}, response.ErrEmailNotFound
	}

	if body.Password != existingUser.Password {
		return models.User{}, response.ErrInvalidCredentials
	}

	return existingUser, nil
}
