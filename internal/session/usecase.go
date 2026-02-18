package session

import (
	"context"
	"time"

	"github.com/aclgo/grpc-jwt/internal/session/models"
	"github.com/golang-jwt/jwt"
)

const (
	TypeAccessTTK    = "access"
	TypeRefreshTTK   = "refresh"
	TtlExpAccessTTK  = time.Minute * 1
	TtlExpRefreshTTK = time.Hour * 24
)

type SessionUC interface {
	CreateTokens(context.Context, string, string) (*models.Token, error)
	RefreshToken(context.Context, string, string) (*models.Token, error)
	ValidToken(context.Context, string) (jwt.MapClaims, error)
	RevogeToken(context.Context, string, string) error
	VerifyRevogedTokens(context.Context, string) error
	GetClaimsRefreshToken(context.Context, string) (jwt.MapClaims, error)
}
