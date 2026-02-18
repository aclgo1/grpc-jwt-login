package session

import "errors"

var (
	ErrTokenRevoged     = errors.New("token revoged")
	ErrTokenExpired     = errors.New("Token is expired")
	ErrInvalidToken     = errors.New("token invalid")
	ErrTypeTokenInvalid = errors.New("type token invalid")
	ErrUnknown          = errors.New("unknown error")
	ErrMistachTokenID   = errors.New("mistach token id")
)
