package usecase

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/aclgo/grpc-jwt/config"
	"github.com/aclgo/grpc-jwt/pkg/logger"
	"github.com/alicebob/miniredis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func SetupRedisTest() *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

func TestCreateTTK(t *testing.T) {
	t.Parallel()

	logger := logger.NewapiLogger(nil)
	redisClient := SetupRedisTest()
	fakeUserID := uuid.NewString()
	fakeRole := "fake-admin"
	sessUC := NewSessionUC(logger, redisClient, "secrey")

	t.Run("Create tokens", func(t *testing.T) {

		tokens, err := sessUC.CreateTokens(context.Background(), fakeUserID, fakeRole)
		require.NoError(t, err)
		require.NotNil(t, tokens)
	})
}

func TestValidTTK(t *testing.T) {
	t.Parallel()

	logger := logger.NewapiLogger(nil)
	redisClient := SetupRedisTest()
	fakeUserID := uuid.NewString()
	fakeRole := "fake-admin"
	sessUC := NewSessionUC(logger, redisClient, "secrey")

	t.Run("Valid tokens", func(t *testing.T) {
		tokens, err := sessUC.CreateTokens(context.Background(), fakeUserID, fakeRole)
		require.NoError(t, err)
		require.NotNil(t, tokens)

		claims, err := sessUC.ValidToken(context.Background(), tokens.Access)
		require.NoError(t, err)
		require.NotNil(t, claims)
	})
}

func TestRefreshTTK(t *testing.T) {
	t.Parallel()

	logger := logger.NewapiLogger(nil)
	redisClient := SetupRedisTest()
	fakeUserID := uuid.NewString()
	fakeRole := "fake-admin"
	sessUC := NewSessionUC(logger, redisClient, "secrey")
	t.Run("Refresh tokens", func(t *testing.T) {
		tokens, err := sessUC.CreateTokens(context.Background(), fakeUserID, fakeRole)
		require.NoError(t, err)
		require.NotNil(t, tokens)

		refresh, err := sessUC.RefreshToken(context.Background(), tokens.Access, tokens.Refresh)
		require.NoError(t, err)
		require.NotNil(t, refresh)
	})
}

func TestRevogeTTK(t *testing.T) {
	t.Parallel()

	logger := logger.NewapiLogger(nil)
	redisClient := SetupRedisTest()
	fakeUserID := uuid.NewString()
	fakeRole := "fake-admin"
	sessUC := NewSessionUC(logger, redisClient, "secrey")
	t.Run("Revoge tokens", func(t *testing.T) {
		tokens, err := sessUC.CreateTokens(context.Background(), fakeUserID, fakeRole)
		require.NoError(t, err)
		require.NotNil(t, tokens)

		errAccess := sessUC.RevogeToken(context.Background(), tokens.Access, tokens.Refresh)
		require.NoError(t, errAccess)

		errRefresh := sessUC.RevogeToken(context.Background(), tokens.Access, tokens.Refresh)
		require.NoError(t, errRefresh)
	})
}

func TestRevogedTTK(t *testing.T) {

	t.Parallel()

	logger := logger.NewapiLogger(&config.Config{})
	redisClient := SetupRedisTest()
	fakeUserID := uuid.NewString()
	fakeRole := "fake-admin"

	sessUC := NewSessionUC(logger, redisClient, "secrey")
	t.Run("Verify if Revoged tokens", func(t *testing.T) {

		tokens, err := sessUC.CreateTokens(context.Background(), fakeUserID, fakeRole)
		require.NoError(t, err)
		require.NotNil(t, tokens)

		errAccess := sessUC.RevogeToken(context.Background(), tokens.Access, tokens.Refresh)
		require.NoError(t, errAccess)

		errRefresh := sessUC.RevogeToken(context.Background(), tokens.Access, tokens.Refresh)
		require.NoError(t, errRefresh)

		err = sessUC.VerifyRevogedTokens(context.Background(), tokens.Access)
		require.Error(t, err, errors.New("token revoged"))
	})
}
