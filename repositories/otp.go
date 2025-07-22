package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/redis/go-redis/v9"
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

// Redis
type redisOTPRepository struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisOTPRepository(client *redis.Client, expiration time.Duration) OTPRepository {
	return &redisOTPRepository{
		client:     client,
		expiration: expiration,
	}
}

func (r *redisOTPRepository) SaveOTP(entry models.OTPEntry) error {
	ctx := context.Background()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("otp:%s", entry.Email)
	return r.client.Set(ctx, key, data, r.expiration).Err()
}

func (r *redisOTPRepository) FindByEmail(email string) (models.OTPEntry, bool) {
	ctx := context.Background()
	
	key := fmt.Sprintf("otp:%s", email)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return models.OTPEntry{}, false
	}

	var entry models.OTPEntry
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return models.OTPEntry{}, false
	}
	return entry, true
}

func (r *redisOTPRepository) VerifyOTP(email, otp string) bool {
	entry, found := r.FindByEmail(email)
	if !found || entry.OTP != otp || time.Now().After(entry.ExpiresAt) {
		return false
	}

	ctx := context.Background()
	key := fmt.Sprintf("otp:%s", email)
	r.client.Del(ctx, key)
	return true
}
