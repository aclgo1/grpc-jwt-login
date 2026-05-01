package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
)

type UserRepoDatabase interface {
	Add(context.Context, *models.User) (*models.User, error)
	FindByID(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
	Update(context.Context, *models.User) (*models.User, error)
	Delete(context.Context, string) error
}

type UserRepoCache interface {
	Set(context.Context, *models.User) error
	Get(context.Context, string) (*models.User, error)
	Del(context.Context, string) error
}

type UserUC interface {
	Register(context.Context, *ParamsCreateUser) (*ParamsOutputUser, error)
	FindByID(context.Context, string) (*ParamsOutputUser, error)
	FindByEmail(context.Context, string) (*ParamsOutputUser, error)
	Update(context.Context, *ParamsUpdateUser) (*ParamsOutputUser, error)
	Delete(context.Context, *ParamsDeleteUser) error
	Login(context.Context, string, string) (*models.Tokens, error)
	Logout(context.Context, *ParamLogoutInput) error
	ValidToken(context.Context, *ParamsValidToken) (*ParamsJwtData, error)
	RefreshTokens(context.Context, *ParamsRefreshTokens) (*RefreshTokens, error)
	GetConnsOnlineUser(ctx context.Context) (int, error)
}

func FormatActiveSessionAccess(s string) string {
	return fmt.Sprintf("active-access-session:%s", s)
}

func FormatActiveSessionRefresh(s string) string {
	return fmt.Sprintf("active-refresh-session:%s", s)
}

func FormatTokenDisconnectChannel(userId string) string {
	now := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf("%s|%s", userId, now)
}

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrEmailCadastred    = errors.New("email cadastred")
	ErrInvalidEmail      = errors.New("email invalid")
	ErrPasswordSmall     = errors.New("password small lenght")
)
