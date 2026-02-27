package usecase

import (
	"context"

	"fmt"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
	session "github.com/aclgo/grpc-jwt/internal/session"
	"github.com/aclgo/grpc-jwt/internal/user"
	"github.com/aclgo/grpc-jwt/internal/utils"
	"github.com/aclgo/grpc-jwt/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrEmailCadastred    = errors.New("email cadastred")
	ErrInvalidEmail      = errors.New("email invalid")
	ErrPasswordSmall     = errors.New("password small lenght")
)

type userUC struct {
	logger           logger.Logger
	userRepoDatabase user.UserRepoDatabase
	userRepoCache    user.UserRepoCache
	jwtSession       session.SessionUC
	rc               *redis.Client
}

func NewUserUC(logger logger.Logger,
	userRepoDatabase user.UserRepoDatabase,
	userRepoCache user.UserRepoCache, sessionUC session.SessionUC, rc *redis.Client) *userUC {
	return &userUC{
		logger:           logger,
		userRepoDatabase: userRepoDatabase,
		userRepoCache:    userRepoCache,
		jwtSession:       sessionUC,
		rc:               rc,
	}
}

func (u *userUC) Register(ctx context.Context, params *user.ParamsCreateUser) (*user.ParamsOutputUser, error) {
	if !utils.ValidMail(params.Email) {
		u.logger.Errorf("Register.FindByEmail: %v", ErrInvalidEmail)
		return nil, fmt.Errorf("Register.FindByEmail: %v", ErrInvalidEmail)
	}

	foundUser, _ := u.userRepoDatabase.FindByEmail(ctx, params.Email)
	if foundUser != nil {
		u.logger.Errorf("Register.FindByEmail: %v", ErrEmailCadastred)
		return nil, fmt.Errorf("Register.FindByEmail: %v", ErrEmailCadastred)
	}

	created, err := u.userRepoDatabase.Add(ctx, &models.User{
		UserID:    uuid.NewString(),
		Name:      params.Name,
		Lastname:  params.Lastname,
		Password:  params.HashPass(),
		Email:     params.Email,
		Role:      string(user.ClientRole),
		Verified:  user.DefaultVerifiedNo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		u.logger.Errorf("Register.Add: %v", err)
		return nil, fmt.Errorf("Register.Add: %v", err)
	}

	return user.Dto(created), nil
}

func (u *userUC) Login(ctx context.Context, email string, password string) (*models.Tokens, error) {

	foundUser, err := u.userRepoDatabase.FindByEmail(ctx, email)
	if err != nil {
		u.logger.Errorf("Login.FindByEmail: %v", err)
		return nil, fmt.Errorf("Login.FindByEmail: %v", err)
	}

	if err := foundUser.ComparePass(password); err != nil {
		u.logger.Errorf("Login: %v", ErrPasswordIncorrect)
		return nil, ErrPasswordIncorrect
	}

	if foundUser.Verified == user.DefaultVerifiedNo {
		return nil, user.ErrUserNotVerified{}
	}

	tokens, err := u.jwtSession.CreateTokens(ctx, foundUser.UserID, foundUser.Role)
	if err != nil {
		u.logger.Errorf("Login.CreateTokens: %v", err)
		return nil, fmt.Errorf("Login.CreateTokens: %v", err)
	}

	pipe := u.rc.Pipeline()

	pipe.Set(ctx, user.FormatActiveSessionAccess(foundUser.UserID), tokens.Access, session.TtlExpAccessTTK)
	pipe.Set(ctx, user.FormatActiveSessionRefresh(foundUser.UserID), tokens.Refresh, session.TtlExpRefreshTTK)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipe.Exec: %w", err)
	}

	return &models.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil

}

func (u *userUC) Logout(ctx context.Context, in *user.ParamLogoutInput) error {
	mc, err := u.jwtSession.ValidToken(ctx, in.AccessToken)
	id, ok := mc["id"].(string)
	if !ok {

	}

	pipe := u.rc.Pipeline()
	pipe.Del(ctx, user.FormatActiveSessionAccess(id))
	pipe.Del(ctx, user.FormatActiveSessionRefresh(id))

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipe.Exec: %w", err)
	}

	err = u.jwtSession.RevogeToken(ctx, in.AccessToken, in.RefreshToken)
	if err != nil {
		u.logger.Errorf("Logout.RevogeToken: %v", err)
		return fmt.Errorf("Logout.RevogeToken: %v", err)
	}

	return nil
}

func (u *userUC) FindByID(ctx context.Context, userID string) (*user.ParamsOutputUser, error) {

	foundUser, err := u.userRepoDatabase.FindByID(ctx, userID)
	if err != nil {
		u.logger.Errorf("FindByID: %v", err)
		return nil, fmt.Errorf("FindByID: %v", err)
	}

	return user.Dto(foundUser), nil
}

func (u *userUC) FindByEmail(ctx context.Context, userEmail string) (*user.ParamsOutputUser, error) {

	// var (
	// 	foundUser *models.User
	// 	err       error
	// )

	// foundUser, err := u.userRepoCache.Get(ctx, userEmail)
	// if err == redis.Nil {
	foundUser, err := u.userRepoDatabase.FindByEmail(ctx, userEmail)
	if err != nil {
		u.logger.Errorf("FindByEmail: %v", err)
		return nil, fmt.Errorf("FindByEmail: %v", err)
	}

	// if err := u.userRepoCache.Set(ctx, foundUser); err != nil {
	// 	u.logger.Warn("FindByEmail.Set: %v", err)
	// }

	return user.Dto(foundUser), nil
	// }

	// if err != nil {
	// 	u.logger.Errorf("FindByEmail.Get: %v", err)
	// 	return nil, fmt.Errorf("FindByEmail.Get: %v", err)
	// }

	// return user.Dto(foundUser), nil
}

func (u *userUC) Update(ctx context.Context, params *user.ParamsUpdateUser) (*user.ParamsOutputUser, error) {

	newUser, err := u.userRepoDatabase.Update(ctx,
		&models.User{
			UserID:    params.UserID,
			Name:      params.Name,
			Lastname:  params.Lastname,
			Password:  params.HashPass(),
			Email:     params.Email,
			Verified:  params.Verified,
			Role:      params.Role,
			UpdatedAt: time.Now(),
		},
	)

	if err != nil {
		u.logger.Errorf("Update.Update: %v", err)
		return nil, errors.Wrap(err, "Update.Update")
	}

	return user.Dto(newUser), nil
}

func (u *userUC) Delete(ctx context.Context, params *user.ParamsDeleteUser) error {
	return u.userRepoDatabase.Delete(ctx, params.UserID)
}

func (u *userUC) ValidToken(ctx context.Context, params *user.ParamsValidToken) (*user.ParamsJwtData, error) {
	claims, err := u.jwtSession.ValidToken(ctx, params.AccessToken)
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("ValidToken: %v", err)
		return nil, err
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return nil, user.ErrInvalidTokenClaims{}
	}

	activeSession, err := u.rc.Get(ctx, user.FormatActiveSessionAccess(userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, user.ErrSessionExpiredOrLoginNewDisp{}
		}
		return nil, err
	}

	if activeSession != params.AccessToken {
		return nil, user.ErrSessionExpiredOrLoginNewDisp{}
	}

	return &user.ParamsJwtData{
		UserID: userID,
		Role:   claims["role"].(string),
	}, nil
}

func (u *userUC) RefreshTokens(ctx context.Context, params *user.ParamsRefreshTokens) (*user.RefreshTokens, error) {

	mc, err := u.jwtSession.GetClaimsRefreshToken(ctx, params.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("u.jwtSession.GetClaimsRefreshToken: %w", err)
	}
	userID, ok := mc["id"].(string)

	if !ok {
		return nil, errors.New("failed get user id for refresh token")
	}

	ractive, err := u.rc.Get(ctx, user.FormatActiveSessionRefresh(userID)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err == redis.Nil || ractive != params.RefreshToken {
		return nil, user.ErrSessionExpiredOrLoginNewDisp{}
	}

	tokens, err := u.jwtSession.RefreshToken(ctx, params.AccessToken, params.RefreshToken)
	if err != nil {
		return nil, err
	}

	pipe := u.rc.Pipeline()

	pipe.Set(ctx, user.FormatActiveSessionAccess(userID), tokens.Access, session.TtlExpAccessTTK)
	pipe.Set(ctx, user.FormatActiveSessionRefresh(userID), tokens.Refresh, session.TtlExpRefreshTTK)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipe.Exec: %w", err)
	}

	out := user.RefreshTokens{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}

	return &out, nil
}
