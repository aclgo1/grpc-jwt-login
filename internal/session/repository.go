package session

import (
	"context"
	"fmt"
	"time"
)

type TokenRepo interface {
	Get(ctx context.Context, tokenString string) error
	Set(ctx context.Context, tokenString string, ttl time.Duration) error
}

func FormatAccessTokenRepo(tk string) string {
	return fmt.Sprintf("access-token:%s", tk)
}

func FormatRefreshTokenRepo(tk string) string {
	return fmt.Sprintf("refresh-token:%s", tk)
}

func FormatAccessTokenSessionRepo(tk string) string {
	return fmt.Sprintf("session-access-token:%s", tk)
}

func FormatRefreshTokenSessionRepo(tk string) string {
	return fmt.Sprintf("session-refresh-token:%s", tk)
}
