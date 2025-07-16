package repositories

import (
	"sync"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/models"
)

type OTPRepository interface {
	FindByEmail(email string) (models.OTPEntry, bool)
	SaveOTP(otpEntry models.OTPEntry) error
	VerifyOTP(email string, otp string) bool
}

// In Memory
type inMemoryOTPRepository struct {
	store map[string]models.OTPEntry
	mu sync.RWMutex
}

func NewInMemoryOTPRepository() OTPRepository {
	return &inMemoryOTPRepository{
		store: make(map[string]models.OTPEntry),
	}
}

func (r *inMemoryOTPRepository) FindByEmail(email string) (models.OTPEntry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.store[email]
	return entry, exists
}

func (r *inMemoryOTPRepository) SaveOTP(entry models.OTPEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[entry.Email] = entry
	return nil
}

func (r *inMemoryOTPRepository) VerifyOTP(email, otp string) (bool) {
	r.mu.RLock()
	entry, exists := r.store[email]
	r.mu.RUnlock()

	if !exists {
		return false
	}
	if time.Now().After(entry.ExpiresAt) {
		r.mu.Lock()
		delete(r.store, email)
		r.mu.Unlock()
		return false
	}
	if entry.OTP != otp {
		return false
	}

	r.mu.Lock()
	delete(r.store, email)
	r.mu.Unlock()
	return true
}
