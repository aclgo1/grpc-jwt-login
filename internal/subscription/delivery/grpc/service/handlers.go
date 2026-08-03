package service

import (
	"github.com/aclgo/grpc-jwt/internal/subscription"
	"github.com/aclgo/grpc-jwt/proto"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func(s *subscriptionService) CreateOrExtend(ctx context.Context, req *proto.CreateOrExtendSubscriptionRequest)(*proto.CreateOrExtendSubscriptionResponse,error) {
	pu := subscription.SubscriptionInput{
		UserId: req.GetUserId(),
		Plan: req.GetPlan(),
		Days: req.GetDays(),
	}

	if err := pu.Validate();err != nil {
		return nil,err
	}


	ce, err := s.subscriptionUC.CreateOrExtend(ctx, &pu)
	if err != nil {
		return nil, err
	}

	out := proto.CreateOrExtendSubscriptionResponse{
    	Id:ce.Id,
    	UserId:   ce.UserId,
    	Plan:  ce.Plan,
    	Status:    ce.Status,
    	StartsAt:   timestamppb.New(ce.StartAt),
    	ExpiresAt: timestamppb.New(ce.ExpiresAt),
    	CreatedAt:timestamppb.New(ce.CreatedAt),
    	UpdatedAt:timestamppb.New(ce.UpdatedAt),
	}

	return &out,nil
}
func(s *subscriptionService) CancelSubscription(ctx context.Context, req *proto.CancelSubscriptionRequest)(*proto.CancelSubscriptionResponse,error){
	pu := subscription.SubscriptionCancelInput{
		UserId: req.GetUserId(),
	}

	if err := pu.Validate(); err != nil {
		return nil, err
	}

	canceled,err := s.subscriptionUC.CancelSubscription(ctx, &pu)
	if err != nil {
		return nil, err
	}

	out := proto.CancelSubscriptionResponse{
		SubscriptionId: canceled.SubscriptionId,
		UserId: canceled.UserId,
		Status: canceled.Status,
		UpdatedAt: timestamppb.New(canceled.UpdatedAt),
		
	}

	return &out,nil
}
func(s *subscriptionService) CheckIsPremium(ctx context.Context, req *proto.CheckIsPremiumRequest)(*proto.CheckIsPremiumResponse,error){
	pu := subscription.SubscriptionIsPremiunInput{
		UserId: req.GetUserId(),
	}

	if err := pu.Validate();err != nil {
		return nil,err
	}

	chk,err := s.subscriptionUC.CheckIsPremiun(ctx,&pu)
	if err != nil {
		return nil,err
	}

	out := proto.CheckIsPremiumResponse{
		Active: chk.IsPremiun,
		ExpiresAt: timestamppb.New(chk.ExpiresAt),
	}

	return &out,nil
}