package subscription

import (
	"context"
	"errors"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)


type SubscriptionUseCase interface{
	CreateOrExtend(context.Context, *SubscriptionInput)(*SubscriptionOutput,error)
	CheckIsPremiun(context.Context, *SubscriptionIsPremiunInput)(*SubscriptionIsPremiunOutput,error)
	CancelSubscription(context.Context, *SubscriptionCancelInput)(*SubscriptionCancelOutput,error)
}

type SubscriptionRepository interface{
	SubscribeExtend(context.Context, *models.SubscriptionInput)(*models.SubscriptionOutput,error)
	CheckPremiun(context.Context, *models.SubscriptionIsPremiunInput)(*models.SubscriptionIsPremiunOutput,error)
	UpdateSubscriptionsStatusExpired(context.Context,*models.SubscriptionExpiredInput) (*models.SubscriptionExpiredOutput,error)
	UpdateSubscriptionsStatusCancel(context.Context,*models.SubscriptionCancelInput) (*models.SubscriptionCancelOutput,error)
}

type SubscriptionRepositoryRedis interface{
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}


type SubscriptionInput struct {
	Id string
	UserId string
	Plan string
	Days int64
}

func(p *SubscriptionInput)Validate()error{

	if p.Id == "" {
		p.Id = uuid.NewString()
	}

	_, err := uuid.Parse(p.Id)
	if err != nil {
		return err
	}

	if p.UserId == ""{
		return errors.New("user id empty")
	}

	var plan string
	var days int64
	switch p.Plan{
	case "7_days":
		plan = p.Plan
		days = 7
	case "1_month":
		plan = p.Plan
		days = 30
	case "1_year":
		plan = p.Plan
		days = 365
	case "undefined":
		plan = p.Plan
		if p.Days <= 0 {
			return errors.New("plan days invalid")
		}

		days = p.Days

	default:
		return errors.New("plan subscription invalid")
	}

	p.Plan = plan
	p.Days = days

	return nil
}

type SubscriptionOutput struct {
	Id string 
	UserId string  
	Plan string 
	Status string 
	StartAt time.Time 
	ExpiresAt time.Time 
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

type SubscriptionIsPremiunInput struct{
	UserId string
}

func(p *SubscriptionIsPremiunInput)Validate()error{
	return nil
}

type SubscriptionIsPremiunOutput struct{
	IsPremiun bool
	ExpiresAt time.Time
}

type SubscriptionCancelInput struct{
	UserId string
}

func(p *SubscriptionCancelInput)Validate()error{
	return nil
}

type SubscriptionCancelOutput struct{
	SubscriptionId string    
	UserId         string    
	Status         string    
	UpdatedAt      time.Time
}