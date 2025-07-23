package services

import (
	"time"

	"github.com/MegeKaplan/megebase-identity-service/dto"
	"github.com/MegeKaplan/megebase-identity-service/messaging"
	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/MegeKaplan/megebase-identity-service/repositories"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(body dto.RegisterRequest) (models.User, *response.AppError)
	LoginUser(body dto.LoginRequest) (models.User, *response.AppError)
	SendOTP(body dto.SendOTPRequest) *response.AppError
	VerifyOTP(email string, otp string) (bool, *response.AppError)
	GenerateTokens(user models.User) (string, string, *response.AppError)
	RefreshTokens(refreshToken string) (string, string, *response.AppError)
}

type authService struct {
	authRepo         repositories.UserRepository
	otpRepo          repositories.OTPRepository
	messagingService messaging.MessagingService
	refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService(authRepo repositories.UserRepository, otpRepo repositories.OTPRepository, messagingService messaging.MessagingService, refreshTokenRepo repositories.RefreshTokenRepository) AuthService {
	return &authService{
		authRepo:         authRepo,
		otpRepo:          otpRepo,
		messagingService: messagingService,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *authService) RegisterUser(body dto.RegisterRequest) (models.User, *response.AppError) {
	_, err := s.authRepo.FindByEmail(body.Email)
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

	if err := s.authRepo.Create(&user); err != nil {
		return models.User{}, response.ErrDBOperation
	}

	return user, nil
}

func (s *authService) LoginUser(body dto.LoginRequest) (models.User, *response.AppError) {
	existingUser, err := s.authRepo.FindByEmail(body.Email)
	if err != nil {
		return models.User{}, response.ErrEmailNotFound
	}

	if !utils.CheckPasswordHash(body.Password, existingUser.Password) {
		return models.User{}, response.ErrInvalidCredentials
	}

	return existingUser, nil
}

func (s *authService) SendOTP(body dto.SendOTPRequest) *response.AppError {
	entry, exists := s.otpRepo.FindByEmail(body.Email)
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

	err = s.otpRepo.SaveOTP(otpEntry)
	if err != nil {
		return response.ErrDBOperation
	}

	if err := utils.SendOTP(s.messagingService, "email", body.Email, otp); err != nil {
		return response.ErrOTPSendFailed
	}

	return nil
}

func (s *authService) VerifyOTP(email string, otp string) (bool, *response.AppError) {
	entry, exists := s.otpRepo.FindByEmail(email)
	if !exists {
		return false, response.ErrOTPNotFound
	}

	if time.Now().After(entry.ExpiresAt) {
		return false, response.ErrOTPExpired
	}

	if entry.OTP != otp {
		return false, response.ErrInvalidOTP
	}

	isValid := s.otpRepo.VerifyOTP(email, otp)

	return isValid, nil
}

func (s *authService) GenerateTokens(user models.User) (string, string, *response.AppError) {
	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", "", response.ErrTokenGenerationFailed
	}

	refreshToken := utils.GenerateRefreshToken()

	refreshEntry := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID.String(),
		ExpiresAt: time.Now().Add(utils.RefreshTokenTTL()),
	}

	if err := s.refreshTokenRepo.Save(refreshEntry); err != nil {
		return "", "", response.ErrDBOperation
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshTokens(refreshToken string) (string, string, *response.AppError) {
	entry, found := s.refreshTokenRepo.Find(refreshToken)
	if !found {
		return "", "", response.ErrInvalidRefreshToken
	}

	if time.Now().After(entry.ExpiresAt) {
		return "", "", response.ErrExpiredRefreshToken
	}

	if err := s.refreshTokenRepo.Delete(refreshToken); err != nil {
		return "", "", response.ErrDBOperation
	}

	userID, err := uuid.Parse(entry.UserID)
	if err != nil {
		return "", "", response.ErrInvalidUserID
	}

	user, err := s.authRepo.FindByID(userID.String())
	if err != nil {
		return "", "", err.(*response.AppError)
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err.(*response.AppError) != nil {
		return "", "", response.ErrTokenGenerationFailed
	}

	newRefreshToken := utils.GenerateRefreshToken()

	refreshEntry := models.RefreshToken{
		Token:     newRefreshToken,
		UserID:    user.ID.String(),
		ExpiresAt: time.Now().Add(utils.RefreshTokenTTL()),
	}

	if err := s.refreshTokenRepo.Save(refreshEntry); err != nil {
		return "", "", response.ErrDBOperation
	}

	return accessToken, newRefreshToken, nil
}
