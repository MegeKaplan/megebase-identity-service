package services

import (
	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
)

type UserService interface {
	GetUserByID(id string) (models.User, *response.AppError)
	SearchUsers(params utils.QueryParams) ([]models.User, *response.AppError)
	UpdateUser(id string, body dto.UpdateUserRequest) (models.User, *response.AppError)
	DeleteUser(id string, hardDelete bool) *response.AppError
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

func (s *userService) SearchUsers(params utils.QueryParams) ([]models.User, *response.AppError) {
	users, err := s.userRepo.SearchUsers(params)
	if err != nil {
		return nil, response.ErrUsersNotFound
	}
	return users, nil
}

func (s *userService) UpdateUser(id string, body dto.UpdateUserRequest) (models.User, *response.AppError) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return models.User{}, response.ErrUserNotFound
	}

	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.Password != nil {
		hashedPassword, err := utils.HashPassword(*body.Password)
		if err != nil {
			return models.User{}, err
		}
		user.Password = hashedPassword
	}
	if body.Name != nil {
		user.Name = *body.Name
	}

	if err := s.userRepo.Update(user); err != nil {
		return models.User{}, response.ErrUserUpdateFailed
	}

	return user, nil
}

func (s *userService) DeleteUser(id string, hardDelete bool) *response.AppError {
	if hardDelete {
		if err := s.userRepo.HardDeleteByID(id); err != nil {
			return response.ErrUserDeleteFailed
		}
	} else {
		if err := s.userRepo.SoftDeleteByID(id); err != nil {
			return response.ErrUserDeleteFailed
		}
	}
	return nil
}
