package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	jwtSecret      []byte
	jwtExpireHours int
	isInitialized bool = false
)

func ConfigureJWT() {
	if isInitialized {
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set in .env")
	}

	jwtSecret = []byte(secret)

	expireHoursStr := os.Getenv("JWT_EXPIRE_HOURS")
	if expireHoursStr == "" {
		expireHoursStr = "24"
	}

	expireHours, err := strconv.Atoi(expireHoursStr)
	if err != nil {
		panic("JWT_EXPIRE_HOURS must be an integer")
	}

	jwtExpireHours = expireHours

	isInitialized = true
}

type JWTClaim struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, email string) (string, *response.AppError) {
	ConfigureJWT()
	claims := JWTClaim{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "megebase-identity-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", response.ErrTokenGenerationFailed
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (*JWTClaim, *response.AppError) {
	ConfigureJWT()
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, response.ErrInvalidSigningMethod
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, response.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, response.ErrInvalidToken
	}

	return claims, nil
}
