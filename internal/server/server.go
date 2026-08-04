package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/aclgo/grpc-jwt/config"
	"github.com/aclgo/grpc-jwt/internal/interceptor"
	sessionUC "github.com/aclgo/grpc-jwt/internal/session/usecase"
	subscriptionService "github.com/aclgo/grpc-jwt/internal/subscription/delivery/grpc/service"
	subscriptionRepo "github.com/aclgo/grpc-jwt/internal/subscription/repository"
	subscriptionUC "github.com/aclgo/grpc-jwt/internal/subscription/usecase"
	"github.com/aclgo/grpc-jwt/internal/user/delivery/grpc/service"
	userRepo "github.com/aclgo/grpc-jwt/internal/user/repository"
	userUC "github.com/aclgo/grpc-jwt/internal/user/usecase"
	"github.com/aclgo/grpc-jwt/pkg/grpc_auth"
	"github.com/aclgo/grpc-jwt/pkg/logger"
	"github.com/aclgo/grpc-jwt/proto"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server struct {
	db          *sqlx.DB
	redisClient *redis.Client
	logger      logger.Logger
	config      *config.Config
}

func NewServer(db *sqlx.DB, redisClient *redis.Client,
	logger logger.Logger, config *config.Config) *Server {
	return &Server{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		config:      config,
	}
}

func (s *Server) Run() error {
	interceptor := interceptor.NewInterceptor(s.logger)

	sessUC := sessionUC.NewSessionUC(s.logger, s.redisClient, s.config.SecretKey)

	usRepo := userRepo.NewPostgresRepo(s.db)
	usRepoRedis := userRepo.NewredisRepo(s.redisClient)
	userUC := userUC.NewUserUC(s.logger, usRepo, usRepoRedis, sessUC, s.redisClient)
	
	subRepoRedis := subscriptionRepo.NewSubscriptionRedisRepo(s.redisClient)
	subRepo := subscriptionRepo.NewSubscriptionRepo(s.db)
	subUC := subscriptionUC.NewSubscriptionUseCase(context.Background(), 500, time.Hour*12, subRepo,subRepoRedis)


	userService := service.NewUserService(s.logger, userUC)
	subService := subscriptionService.NewSubscriptionService(subUC)

	auth := grpc_auth.NewGrpcAuth(s.config)

	listen, err := net.Listen("tcp", ":"+s.config.ServerPort)

	if err != nil {
		s.logger.Errorf("net.Listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Logger),
		grpc.ChainUnaryInterceptor(auth.AuthInterceptor),
	}

	server := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(server, userService)
	proto.RegisterSubscriptionServiceServer(server,subService)
	s.logger.Infof("server starting port %s", s.config.ServerPort)
	if err := server.Serve(listen); err != nil {
		return fmt.Errorf("Run.NewServer: %v", err)
	}

	return nil
}
