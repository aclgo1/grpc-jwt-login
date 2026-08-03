package service

import (
	"github.com/aclgo/grpc-jwt/internal/subscription"
	"github.com/aclgo/grpc-jwt/proto"
)


type subscriptionService struct {
	subscriptionUC subscription.SubscriptionUseCase
	proto.UnimplementedSubscriptionServiceServer
}

func NewSubscriptionService(sub subscription.SubscriptionUseCase)*subscriptionService{
	return &subscriptionService{
		subscriptionUC: sub,
	}
}