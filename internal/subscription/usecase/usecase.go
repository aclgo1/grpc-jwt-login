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
		Id: input.Id,
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

const unlockLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

func (u *subscriptionUC) expiredSubscriptionBatch(ctx context.Context) {
	ticker := time.NewTicker(u.cronTime)
	defer ticker.Stop()

	lockKey := "lock:expired_subscriptions_job"

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			u.processExpiredBatchJob(ctx, lockKey)
		}
	}
}

func (u *subscriptionUC) processExpiredBatchJob(ctx context.Context, lockKey string) {
	lockToken := uuid.NewString()

	lockTtl := max(u.cronTime-30*time.Second, 5*time.Second)

	acquired, err := u.repoRedis.SetNX(ctx, lockKey, lockToken, lockTtl).Result()
	if err != nil {
		log.Printf("failed to acquire redis lock: %v", err)
		return
	}
	if !acquired {
		return 
	}

	defer func() {
		u.repoRedis.Eval(context.Background(), unlockLuaScript, []string{lockKey}, lockToken)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			p := models.SubscriptionExpiredInput{
				BatchSize: u.BatchSize,
			}

			resp, err := u.repo.UpdateSubscriptionsStatusExpired(ctx, &p)
			if err != nil {
				log.Printf("failed to update expired subscriptions: %v", err)
				return
			}

			if resp == nil || resp.RowsAffected < u.BatchSize {
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}