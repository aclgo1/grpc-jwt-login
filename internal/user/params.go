package user

import (
	"context"
	"errors"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/aclgo/grpc-jwt/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	ClientRole         = "client"
	DefaultVerifiedNo  = "no"
	DefaultVerifiedYes = "yes"
)

type ParamsCreateUser struct {
	Name     string
	Lastname string
	Password string
	Email    string
}

func (p *ParamsCreateUser) HashPass() string {
	bc, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	return string(bc)
}

func (p *ParamsCreateUser) Validate(ctx context.Context) error {
	return nil
}

type ParamsUpdateUser struct {
	UserID    string
	Name      string
	Lastname  string
	Password  string
	Email     string
	Verified  string
	Role      string
	UpdatedAt time.Time
}

func (p *ParamsUpdateUser) HashPass() string {
	if p.Password == "" {
		return ""
	}

	bc, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	return string(bc)
}

func (p *ParamsUpdateUser) Validate() error {
	if p.UserID == "" {
		return errors.New("user id empty")
	}

	if p.Email != "" {
		if !utils.ValidMail(p.Email) {
			return ErrUserInvalidEmail{}
		}
	}

	return nil
}

type ParamsOutputUser struct {
	Id        string
	Name      string
	Lastname  string
	Password  string
	Email     string
	Role      string
	Verified  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *ParamsOutputUser) ClearPass() {
	p.Password = ""
}

func Dto(user *models.User) *ParamsOutputUser {
	return &ParamsOutputUser{
		Id:        user.UserID,
		Name:      user.Name,
		Lastname:  user.Lastname,
		Password:  "",
		Email:     user.Email,
		Role:      user.Role,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type ParamsJwtData struct {
	UserID string
	Role   string
}

type ParamsValidToken struct {
	AccessToken string
}

type ParamsDeleteUser struct {
	UserID string
}

type ParamsRefreshTokens struct {
	AccessToken  string
	RefreshToken string
}

func (p *ParamsRefreshTokens) Validate() error {
	if p.AccessToken == "" {
		return errors.New("access token empty")
	}

	if p.RefreshToken == "" {
		return errors.New("refresh token empty")
	}

	return nil
}

type RefreshTokens struct {
	AccessToken  string
	RefreshToken string
}

func (p *RefreshTokens) Validate() error {
	if p.AccessToken == "" {
		return errors.New("access token empty")
	}

	if p.RefreshToken == "" {
		return errors.New("refresh token empty")
	}

	return nil
}

type ParamLogoutInput struct {
	AccessToken  string
	RefreshToken string
}

func (p *ParamLogoutInput) Validate() error {
	if p.AccessToken == "" {
		return errors.New("access token empty")
	}

	if p.RefreshToken == "" {
		return errors.New("refresh token empty")
	}

	return nil
}

type ErrUserNotVerified struct {
}

func (e ErrUserNotVerified) Error() string {
	return "user not verified"
}

type ErrUserInvalidEmail struct {
}

func (e ErrUserInvalidEmail) Error() string {
	return "invalid email"
}

type ErrLoginNewDisp struct {
}

func (e ErrLoginNewDisp) Error() string {
	return "login in new dispositivy"
}

type ErrSessionExpired struct {
}

func (e ErrSessionExpired) Error() string {
	return "session expired"
}

type ErrInvalidTokenClaims struct {
}

func (e ErrInvalidTokenClaims) Error() string {
	return "invalid token claims"
}
