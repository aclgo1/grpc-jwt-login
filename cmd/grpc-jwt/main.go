package main

import (
	"log"

	"github.com/aclgo/grpc-jwt/config"
	"github.com/aclgo/grpc-jwt/internal/server"
	"github.com/aclgo/grpc-jwt/internal/session"
	"github.com/aclgo/grpc-jwt/migrations"
	"github.com/aclgo/grpc-jwt/pkg/logger"
	pts "github.com/aclgo/grpc-jwt/pkg/postgres"
	rredis "github.com/aclgo/grpc-jwt/pkg/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(cfg)

	logger := logger.NewapiLogger(cfg)
	logger.InitLogger()
	logger.Info("logger initialized")

	db, err := pts.Connect(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	srcDriver,err := iofs.New(migrations.MigrationsFs, ".")
	if err != nil {
		log.Fatalf("iofs.New: %v\n",err)
	}

	dbDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil{
		log.Fatalf("postgres.NewIstance: %v\n",err)
	}

	m,err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"postgres",
		dbDriver,
	)

	if err != nil {
		log.Fatalf("migrate.NewInstance: %v\n",err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange{
		log.Fatalf("up migrate: %v\n",err)
	}


	redisClient := rredis.NewRedisClient(cfg)

	session.SetSettingsSession(cfg.TimeExpirateAccessToken, cfg.TimeExpirateRefreshToken)

	server := server.NewServer(db, redisClient, logger, cfg)

	logger.Fatal(server.Run())

}
