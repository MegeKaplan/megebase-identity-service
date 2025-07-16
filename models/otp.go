package models

import "time"

type OTPEntry struct {
	OTP       string
	Email     string
	ExpiresAt time.Time
}
