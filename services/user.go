package services

import (
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
)

type UserService interface {
	GetUserByID(id string) (models.User, *response.AppError)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id string) (models.User, *response.AppError) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return models.User{}, response.ErrUserNotFound
	}
	return user, nil
}
