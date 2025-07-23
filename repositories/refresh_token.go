package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MegeKaplan/megebase-identity-service/models"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepository interface {
	Save(entry models.RefreshToken) error
	Find(token string) (models.RefreshToken, bool)
	Delete(token string) error
}

type refreshTokenRepository struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRefreshTokenRepository(client *redis.Client, expiration time.Duration) RefreshTokenRepository {
	return &refreshTokenRepository{
		client:     client,
		expiration: expiration,
	}
}

func (r *refreshTokenRepository) Save(entry models.RefreshToken) error {
	ctx := context.Background()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("refresh:%s", entry.Token)

	return r.client.Set(ctx, key, data, r.expiration).Err()
}

func (r *refreshTokenRepository) Find(token string) (models.RefreshToken, bool) {
	ctx := context.Background()

	key := fmt.Sprintf("refresh:%s", token)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return models.RefreshToken{}, false
	}

	var entry models.RefreshToken
	if err := json.Unmarshal([]byte(val), &entry); err != nil {
		return models.RefreshToken{}, false
	}

	return entry, true
}

func (r *refreshTokenRepository) Delete(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", token)
	return r.client.Del(ctx, key).Err()
}
