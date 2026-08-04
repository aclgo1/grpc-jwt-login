package usecase

import (
	"context"
	"log"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/aclgo/grpc-jwt/internal/subscription"
	"github.com/google/uuid"
)

type subscriptionUC struct {
	BatchSize int64
	cronTime time.Duration
	repo subscription.SubscriptionRepository
	repoRedis subscription.SubscriptionRepositoryRedis
}

func NewSubscriptionUseCase(ctx context.Context, batchSize int64, cronTime time.Duration,
	repo subscription.SubscriptionRepository, repoRedis subscription.SubscriptionRepositoryRedis)subscription.SubscriptionUseCase{
	s := subscriptionUC{
		BatchSize: batchSize,
		cronTime: cronTime,
		repo: repo,
		repoRedis: repoRedis,
	}

	go s.expiredSubscriptionBatch(ctx)

	return &s
}

func(u *subscriptionUC)	CreateOrExtend(ctx context.Context, input *subscription.SubscriptionInput)(*subscription.SubscriptionOutput,error){



	pm := models.SubscriptionInput{
		Id: uuid.NewString(),
		UserId: input.UserId,
		Plan: input.Plan,
		Days: input.Days,
	}

	sub,err := u.repo.SubscribeExtend(ctx, &pm)
	if err != nil {
		return nil, err
	}

	out := subscription.SubscriptionOutput{
		Id: sub.Id,
   		UserId   :sub.UserId ,
   		Plan     :sub.Plan,
   		Status    :sub.Status,
   		StartAt   :sub.StartAt,
   		ExpiresAt :sub.ExpiresAt,
   		CreatedAt :sub.CreatedAt,
   		UpdatedAt :sub.UpdatedAt,
	}

	return &out,nil
}
func(u *subscriptionUC)	CheckIsPremiun(ctx context.Context, input *subscription.SubscriptionIsPremiunInput)(*subscription.SubscriptionIsPremiunOutput,error){

	pm := models.SubscriptionIsPremiunInput{
		UserId: input.UserId,
	}

	chkSub, err := u.repo.CheckPremiun(ctx, &pm)
	if err != nil {
		return nil,err
	}

	out := subscription.SubscriptionIsPremiunOutput{
		IsPremiun: chkSub.Active,
		ExpiresAt: chkSub.ExpiresAt,

	}
	return &out, nil
}
func(u *subscriptionUC)	CancelSubscription(ctx context.Context, input *subscription.SubscriptionCancelInput)(*subscription.SubscriptionCancelOutput,error){

	pm := models.SubscriptionCancelInput{
		UserId: input.UserId,
	}

	cancel, err := u.repo.UpdateSubscriptionsStatusCancel(ctx, &pm)
	if err != nil {
		return nil,err
	}

	out := subscription.SubscriptionCancelOutput{
		SubscriptionId: cancel.SubscriptionId,
		UserId: cancel.UserId,
		Status: cancel.Status,
		UpdatedAt: cancel.UpdatedAt,
	}

	return &out,nil
}

func(u *subscriptionUC)expiredSubscriptionBatch(ctx context.Context){
	ticker := time.NewTicker(u.cronTime)
	defer ticker.Stop()

	lockKey := "lock:expired_subscriptions_job"

	for{
		select{
		case <-ctx.Done():
			return
		case <-ticker.C:

			lockTtl := u.cronTime -time.Second*3
			if lockTtl < time.Second*5{
				lockTtl = time.Second*5
			}

			acquired, err := u.repoRedis.SetNX(ctx,lockKey, "locked", lockTtl).Result()
			if err != nil {
				log.Printf("failed acquired lock redis: %v",err)
				continue
			}

			if !acquired{
				continue
			}

		Loop:
			for{
				select{
				case <-ctx.Done():
					u.repoRedis.Del(ctx, lockKey)
					return
				default:
					p := models.SubscriptionExpiredInput{
						BatchSize: u.BatchSize,
					}

					resp,err := u.repo.UpdateSubscriptionsStatusExpired(ctx, &p)
					if err != nil{
						log.Printf("failed update subscriptions expired: %v\n",err)
						break Loop
					}

					if resp == nil || resp.RowsAffected < u.BatchSize{
						break Loop
					}

					time.Sleep(time.Millisecond *100)
				}
			}
		}
	}
}