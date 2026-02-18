package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aclgo/grpc-jwt/internal/session"
	"github.com/alicebob/miniredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func SetupRedisTest() session.TokenRepo {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return NewjwtStore(client)
}

func TestJwtRepo(t *testing.T) {
	t.Parallel()
	jwtRepo := SetupRedisTest()

	t.Run("Set jwt", func(t *testing.T) {
		err := jwtRepo.Set(context.Background(), "my-token-string", time.Hour)
		require.NoError(t, err)
	})

	t.Run("Get jwt", func(t *testing.T) {
		err := jwtRepo.Get(context.Background(), "my-token-string")
		require.NoError(t, err)
	})
}
