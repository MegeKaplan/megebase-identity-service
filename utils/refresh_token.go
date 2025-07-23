package utils

import (
	"time"

	"github.com/google/uuid"
)

func GenerateRefreshToken() string {
	return uuid.NewString()
}

func RefreshTokenTTL() time.Duration {
	return 7 * 24 * time.Hour
}
