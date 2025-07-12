package utils

import (
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, *response.AppError) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return  "", response.ErrPasswordHashingFailed
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}