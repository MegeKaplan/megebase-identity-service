package services

import (
	"time"

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

func SendOTP(repo repositories.OTPRepository, body dto.SendOTPRequest) *response.AppError {
	entry, exists := repo.FindByEmail(body.Email)
	if exists && time.Now().Before(entry.ExpiresAt) {
		return response.ErrOTPSentRecently
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return response.ErrOTPGenerationFailed
	}

	otpEntry := models.OTPEntry{
		OTP:       otp,
		Email:     body.Email,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	err = repo.SaveOTP(otpEntry)
	if err != nil {
		return response.ErrDBOperation
	}

	if err := utils.SendOTP(otp, body.Email); err != nil {
		return response.ErrOTPSendFailed
	}

	return nil
}

func VerifyOTP(repo repositories.OTPRepository, email string, otp string) (bool, *response.AppError) {
	entry, exists := repo.FindByEmail(email)
	if !exists {
		return false, response.ErrOTPNotFound
	}

	if time.Now().After(entry.ExpiresAt) {
		return false, response.ErrOTPExpired
	}

	if entry.OTP != otp {
		return false, response.ErrInvalidOTP
	}

	isValid := repo.VerifyOTP(email, otp)

	return isValid, nil
}
