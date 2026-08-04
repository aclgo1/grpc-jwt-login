package repository

import (
	"github.com/aclgo/grpc-jwt/internal/subscription"
	"github.com/redis/go-redis/v9"
)


type subscriptionRedisRepo struct{
	*redis.Client
}

func NewSubscriptionRedisRepo(r *redis.Client)subscription.SubscriptionRepositoryRedis{
	return r
}
