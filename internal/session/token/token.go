package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type tokenAction struct {
	SecretKey string
}

func NewtokenAction(secretKey string) *tokenAction {
	return &tokenAction{
		SecretKey: secretKey,
	}
}

func (j *tokenAction) NewToken(typeTTK, userID, role string, ttl time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":      typeTTK,
		"id":        userID,
		"role":      role,
		"createdAt": time.Now().Unix(),
		"exp":       time.Now().Add(ttl).Unix(),
	})

	stoken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", fmt.Errorf("token.SignedString: %v", err)
	}

	return stoken, nil
}

func (j *tokenAction) ParseToken(tokenString string) (*jwt.Token, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = j.parseHS256(tokenString, token)
	if err != nil {
		return nil, fmt.Errorf("parseHS256: %v", err)
	}

	return token, nil
}

func (j *tokenAction) parseHS256(tokenString string, token *jwt.Token) (*jwt.Token, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}

	parse, err := parser.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.SecretKey), nil
	})

	return parse, err
}

func (j *tokenAction) GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if !token.Valid {
		return nil, fmt.Errorf("token.Valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token.Claims.(jwt.MapClaims)")
	}

	return claims, nil
}

func (j *tokenAction) IsExpired(timeUnix float64) bool {
	return float64(time.Now().Unix()) > timeUnix
}
