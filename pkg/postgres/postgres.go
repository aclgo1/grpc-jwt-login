package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/aclgo/grpc-jwt/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open(cfg.DatabaseDriver, cfg.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect %v", err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("PingContext %v", err)
	}

	conn.SetMaxIdleConns(15)
	conn.SetMaxOpenConns(25)
	conn.SetConnMaxLifetime(time.Minute * 5)

	return conn, nil
}
