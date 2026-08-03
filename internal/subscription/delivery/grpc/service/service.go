package service

import "github.com/aclgo/grpc-jwt/internal/subscription"


type subscriptionService struct {
	subscriptionUC subscription.SubscriptionUseCase 
}

func NewSubscriptionService(sub subscription.SubscriptionUseCase)*subscriptionService{
	return &subscriptionService{
		subscriptionUC: sub,
	}
}