package services

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
)

func RegisterUser(repo repositories.UserRepository, body dto.RegisterRequest) (models.User, *response.AppError) {
	_, err := repo.FindByEmail(body.Email)
	if err == nil {
		return models.User{}, response.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err.(*response.AppError) != nil {
		return models.User{}, err.(*response.AppError)
	}

	user := models.User{
		Email:    body.Email,
		Password: hashedPassword,
	}

	if err := repo.Create(&user); err != nil {
		return models.User{}, response.ErrDBOperation
	}

	return user, nil
}

func LoginUser(repo repositories.UserRepository, body dto.LoginRequest) (models.User, *response.AppError) {
	existingUser, err := repo.FindByEmail(body.Email)
	if err != nil {
		return models.User{}, response.ErrEmailNotFound
	}

	if !utils.CheckPasswordHash(body.Password, existingUser.Password) {
		return models.User{}, response.ErrInvalidCredentials
	}

	return existingUser, nil
}
