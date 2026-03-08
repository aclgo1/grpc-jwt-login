package usecase

import (
	"context"
	"fmt"

	session "github.com/aclgo/grpc-jwt/internal/session"
	"github.com/aclgo/grpc-jwt/internal/session/models"
	sessionRepo "github.com/aclgo/grpc-jwt/internal/session/repository"
	sessionToken "github.com/aclgo/grpc-jwt/internal/session/token"
	"github.com/aclgo/grpc-jwt/pkg/logger"

	"github.com/golang-jwt/jwt"
	redis "github.com/redis/go-redis/v9"
)

type sessionUC struct {
	logger      logger.Logger
	tokenRepo   session.TokenRepo
	tokenAction session.TokenAction
}

func NewSessionUC(logger logger.Logger, redisClient *redis.Client,
	secretKey string) *sessionUC {

	tokenRepo := sessionRepo.NewjwtStore(redisClient)
	tokenAction := sessionToken.NewtokenAction(secretKey)

	return &sessionUC{
		logger:      logger,
		tokenRepo:   tokenRepo,
		tokenAction: tokenAction,
	}
}

func (s *sessionUC) CreateTokens(ctx context.Context, userID, role string) (*models.Token, error) {
	return s.createTokens(ctx, userID, role)
}

func (s *sessionUC) RefreshToken(ctx context.Context, accessTTK, refreshTTK string) (*models.Token, error) {

	err := s.verifyRevogedToken(ctx, session.FormatRefreshTokenRepo(refreshTTK))

	if err == nil {
		return nil, session.ErrTokenRevoged
	}

	if err != redis.Nil {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	parsedAccess, err := s.tokenAction.ParseToken(accessTTK)
	if err != nil && err != session.ErrTokenExpired {

		s.logger.Errorf("RefreshToken.ParseToken: %v", err)
		return nil, fmt.Errorf("RefreshToken.ParseToken: %v", err)

	}

	parsedRefresh, err := s.tokenAction.ParseToken(refreshTTK)
	if err != nil {
		s.logger.Errorf("RefreshToken.ParseToken: %v", err)
		return nil, fmt.Errorf("RefreshToken.ParseToken: %v", err)
	}

	claimsAccess, err := s.tokenAction.GetClaims(parsedAccess)
	if err != nil {
		s.logger.Errorf("RefreshToken.GetClaims: %v", err)
		return nil, fmt.Errorf("RefreshToken.GetClaims: %v", err)
	}

	claimsRefresh, err := s.tokenAction.GetClaims(parsedRefresh)
	if err != nil {
		s.logger.Errorf("RefreshToken.GetClaims: %v", err)
		return nil, fmt.Errorf("RefreshToken.GetClaims: %v", err)
	}

	expUnixRefresh := claimsRefresh["exp"].(float64)
	if s.tokenAction.IsExpired(expUnixRefresh) {
		return nil, session.ErrTokenExpired
	}

	idAccess, _ := claimsAccess["id"].(string)
	idRefresh, _ := claimsRefresh["id"].(string)
	typeRefresh, _ := claimsRefresh["type"].(string)

	if idAccess != idRefresh {
		return nil, session.ErrMistachTokenID
	}

	if typeRefresh != session.TypeRefreshTTK {
		return nil, session.ErrTypeTokenInvalid
	}

	role, _ := claimsAccess["role"].(string)

	// if err := s.tokenRepo.Set(ctx, accessTTK, session.TtlExpAccessTTK); err != nil {
	// 	return nil, fmt.Errorf("s.tokenRepo.Set: %w", err)
	// }

	if err := s.tokenRepo.Set(ctx, session.FormatRefreshTokenRepo(refreshTTK), session.TtlExpRefreshTTK); err != nil {
		return nil, fmt.Errorf("s.tokenRepo.Set: %w", err)
	}

	return s.createTokens(ctx, idRefresh, role)

}

func (s *sessionUC) ValidToken(ctx context.Context, ttkString string) (jwt.MapClaims, error) {

	err := s.verifyRevogedToken(ctx, session.FormatAccessTokenRepo(ttkString))
	if err == redis.Nil {
		parsedAccess, err := s.tokenAction.ParseToken(ttkString)
		if err != nil {
			s.logger.Error(session.ErrInvalidToken)
			return nil, session.ErrInvalidToken
		}

		claimsAccess, err := s.tokenAction.GetClaims(parsedAccess)
		if err != nil {
			s.logger.Error(session.ErrInvalidToken)
			return nil, session.ErrInvalidToken
		}

		expUnix := claimsAccess["exp"].(float64)
		if s.tokenAction.IsExpired(expUnix) {
			s.logger.Error(session.ErrTokenExpired)
			return nil, session.ErrTokenExpired
		}

		if claimsAccess["type"].(string) != string(session.TypeAccessTTK) {
			return nil, session.ErrTypeTokenInvalid
		}

		return claimsAccess, nil
	}

	if err == nil {
		return nil, session.ErrTokenRevoged
	}

	return nil, err
}

func (s *sessionUC) GetClaimsRefreshToken(ctx context.Context, ttkString string) (jwt.MapClaims, error) {

	err := s.verifyRevogedToken(ctx, session.FormatRefreshTokenRepo(ttkString))

	if err == nil {
		return nil, session.ErrTokenRevoged
	}

	if err != redis.Nil {
		return nil, err
	}

	parsedRefresh, err := s.tokenAction.ParseToken(ttkString)
	if err != nil {
		s.logger.Error(session.ErrInvalidToken)
		return nil, session.ErrInvalidToken
	}

	claimsRefresh, err := s.tokenAction.GetClaims(parsedRefresh)
	if err != nil {
		s.logger.Error(session.ErrInvalidToken)
		return nil, session.ErrInvalidToken
	}

	expUnix := claimsRefresh["exp"].(float64)
	if s.tokenAction.IsExpired(expUnix) {
		s.logger.Error(session.ErrTokenExpired)
		return nil, session.ErrTokenExpired
	}

	if claimsRefresh["type"].(string) != string(session.TypeRefreshTTK) {
		return nil, session.ErrTypeTokenInvalid
	}

	return claimsRefresh, nil

}

func (s *sessionUC) RevogeToken(ctx context.Context, ttkAccess, ttkRefresh string) error {

	err := s.verifyRevogedToken(ctx, session.FormatAccessTokenRepo(ttkAccess))
	if err == nil {
		return fmt.Errorf("verifyRevogedToken: token access has been revoged")
	}

	err = s.verifyRevogedToken(ctx, session.FormatRefreshTokenRepo(ttkRefresh))
	if err == nil {
		return fmt.Errorf("verifyRevogedToken: token refresh has been revoged")
	}

	errset := s.tokenRepo.Set(
		ctx,
		session.FormatAccessTokenRepo(ttkAccess),
		session.TtlExpAccessTTK,
	)

	if errset != nil {
		s.logger.Errorf("RevogeToken.Set: %v", err)
		return fmt.Errorf("RevogeToken.Set: %v", err)
	}

	errset = s.tokenRepo.Set(ctx,
		session.FormatRefreshTokenRepo(ttkRefresh),
		session.TtlExpRefreshTTK,
	)

	if errset != nil {
		s.logger.Errorf("RevogeToken.Set: %v", err)
		return fmt.Errorf("RevogeToken.Set: %v", err)
	}

	return nil
}

func (s *sessionUC) VerifyRevogedTokens(ctx context.Context, ttkString string) error {
	return s.verifyRevogedToken(ctx, ttkString)
}

func (s *sessionUC) verifyRevogedToken(ctx context.Context, ttkString string) error {
	err := s.tokenRepo.Get(ctx, ttkString)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionUC) createTokens(_ context.Context, userID, role string) (*models.Token, error) {
	access, err := s.tokenAction.NewToken(session.TypeAccessTTK, userID, role, session.TtlExpAccessTTK)
	if err != nil {
		return nil, err
	}

	refresh, err := s.tokenAction.NewToken(session.TypeRefreshTTK, userID, "", session.TtlExpRefreshTTK)
	if err != nil {
		return nil, err
	}

	out := models.Token{
		Access:  access,
		Refresh: refresh,
	}

	return &out, nil
}
