package session

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenAction interface {
	NewToken(typeToken, userID, role string, ttl time.Duration) (string, error)
	ParseToken(tokenString string) (*jwt.Token, error)
	GetClaims(token *jwt.Token) (jwt.MapClaims, error)
	IsExpired(timeUnix float64) bool
}
