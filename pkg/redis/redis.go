package rredis

import (
	"context"
	"log"

	"github.com/aclgo/grpc-jwt/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		DB:       c.RedisDB,
		Password: c.RedisPass,
		PoolSize: 10000,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("NewRedisClient.Ping: %v", err)
	}

	return client
}
