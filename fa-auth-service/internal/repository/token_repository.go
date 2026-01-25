package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiry time.Duration) error
	GetRefreshToken(ctx context.Context, tokenID string) (string, error)
	DeleteRefreshToken(ctx context.Context, tokenID string) error
	DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error
	StorePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error
	GetPasswordResetToken(ctx context.Context, token string) (string, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	// Email verification
	StoreEmailVerificationToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error
	GetEmailVerificationToken(ctx context.Context, token string) (string, error)
	DeleteEmailVerificationToken(ctx context.Context, token string) error
}

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return r.client.Set(ctx, key, userID.String(), expiry).Err()
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return r.client.Get(ctx, key).Result()
}

func (r *tokenRepository) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	key := fmt.Sprintf("refresh:%s", tokenID)
	return r.client.Del(ctx, key).Err()
}

func (r *tokenRepository) DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	pattern := fmt.Sprintf("refresh:*")
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := r.client.Get(ctx, iter.Val()).Result()
		if err == nil && val == userID.String() {
			r.client.Del(ctx, iter.Val())
		}
	}
	return iter.Err()
}

func (r *tokenRepository) StorePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Set(ctx, key, userID.String(), expiry).Err()
}

func (r *tokenRepository) GetPasswordResetToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Get(ctx, key).Result()
}

func (r *tokenRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("reset:%s", token)
	return r.client.Del(ctx, key).Err()
}

// Email verification tokens

func (r *tokenRepository) StoreEmailVerificationToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error {
	key := fmt.Sprintf("verify:%s", token)
	return r.client.Set(ctx, key, userID.String(), expiry).Err()
}

func (r *tokenRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("verify:%s", token)
	return r.client.Get(ctx, key).Result()
}

func (r *tokenRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("verify:%s", token)
	return r.client.Del(ctx, key).Err()
}
