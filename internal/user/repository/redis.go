package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/redis/go-redis/v9"
)

const KeyUser = "user:"

func formatKey(s string) string {
	return fmt.Sprintf("%s%s", KeyUser, s)
}

type redisRepo struct {
	redisClient *redis.Client
}

func NewredisRepo(redisClient *redis.Client) *redisRepo {
	return &redisRepo{
		redisClient: redisClient,
	}
}

func (r *redisRepo) Set(ctx context.Context, user *models.User) error {
	dataUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, formatKey(user.UserID), dataUser, 0).Err()
}

func (r *redisRepo) Get(ctx context.Context, userID string) (*models.User, error) {
	dataUser, err := r.redisClient.Get(ctx, formatKey(userID)).Bytes()
	if err != nil {
		return nil, err
	}

	var user *models.User

	if err := json.Unmarshal(dataUser, &user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *redisRepo) Del(ctx context.Context, userID string) error {
	return r.redisClient.Del(ctx, formatKey(userID)).Err()
}
