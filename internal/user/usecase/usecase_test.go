package usecase

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/aclgo/grpc-jwt/internal/models"
	sessionUC "github.com/aclgo/grpc-jwt/internal/session/usecase"
	"github.com/aclgo/grpc-jwt/internal/user"
	"github.com/aclgo/grpc-jwt/internal/user/mock"
	"github.com/aclgo/grpc-jwt/pkg/logger"
	"github.com/alicebob/miniredis"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func setupRedis() *redis.Client {
	m, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: m.Addr(),
	})

	return client
}

func TestRegister(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepo := mock.NewMockUserPGRepo(ctrl)
	logger := logger.NewapiLogger(nil)

	sessUC := sessionUC.NewSessionUC(logger, setupRedis(), "my-secret-key")

	userUC := NewUserUC(logger, userPGRepo, nil, sessUC, setupRedis())

	uuidUser := uuid.NewString()
	// now := time.Now()

	// mockUser := &models.User{
	// 	Name:     "fake_name",
	// 	Lastname: "fake_lastname",
	// 	Password: "fake_pass",
	// 	Email:    "email@gmail.com",
	// 	Role:     "admin",
	// }

	paramsCreate := user.ParamsCreateUser{
		Name:     "fake_name",
		Lastname: "fake_lastname",
		Password: "fake_pass",
		Email:    "email@gmail.com",
	}

	userPGRepo.EXPECT().FindByEmail(context.Background(), paramsCreate.Email).Return(nil, sql.ErrNoRows)
	userPGRepo.EXPECT().Add(context.Background(), gomock.Any()).Return(&models.User{
		UserID:   uuidUser,
		Name:     "fake_name",
		Lastname: "fake_lastname",
		Password: "fake_pass",
		Email:    "email@gmail.com",
		Role:     "admin",
	}, nil)

	registred, err := userUC.Register(context.Background(), &paramsCreate)
	require.NoError(t, err)
	require.NotNil(t, registred)
	require.Equal(t, uuidUser, registred.Id)
}
func TestFindByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepo := mock.NewMockUserPGRepo(ctrl)
	logger := logger.NewapiLogger(nil)

	sessUC := sessionUC.NewSessionUC(logger, setupRedis(), "my-secret-key")

	userUC := NewUserUC(logger, userPGRepo, nil, sessUC, setupRedis())

	userID := uuid.NewString()

	mockUser := &models.User{
		UserID:   userID,
		Name:     "fake_name",
		Lastname: "fake_lastname",
		Password: "fake_pass",
		Email:    "email@gmail.com",
		Role:     "admin",
	}

	ctx := context.Background()

	userPGRepo.EXPECT().FindByID(ctx, mockUser.UserID).Return(mockUser, nil)

	foundUser, err := userUC.FindByID(ctx, mockUser.UserID)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.Id, mockUser.UserID)

}
func TestFindByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepo := mock.NewMockUserPGRepo(ctrl)
	logger := logger.NewapiLogger(nil)

	sessUC := sessionUC.NewSessionUC(logger, setupRedis(), "my-secret-key")

	userUC := NewUserUC(logger, userPGRepo, nil, sessUC, setupRedis())

	userID := uuid.NewString()

	mockUser := &models.User{
		UserID:   userID,
		Name:     "fake_name",
		Lastname: "fake_lastname",
		Password: "fake_pass",
		Email:    "email@gmail.com",
		Role:     "admin",
	}

	ctx := context.Background()

	userPGRepo.EXPECT().FindByEmail(context.Background(), mockUser.Email).Return(mockUser, nil)

	foundUser, err := userUC.FindByEmail(ctx, mockUser.Email)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.Email, mockUser.Email)

}
func TestUpdate(t *testing.T) {

}
