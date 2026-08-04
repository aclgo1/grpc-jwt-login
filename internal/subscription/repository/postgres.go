package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/aclgo/grpc-jwt/internal/subscription"
	"github.com/jmoiron/sqlx"
)

type subscriptionRepo struct {
	DB *sqlx.DB
}

func NewSubscriptionRepo(db *sqlx.DB) subscription.SubscriptionRepository {
	return &subscriptionRepo{
		DB: db,
	}
}

func (r *subscriptionRepo) SubscribeExtend(ctx context.Context, params *models.SubscriptionInput) (*models.SubscriptionOutput, error) {
	const query = `
        INSERT INTO subscriptions (
            subscription_id, 
            user_id, 
            plan, 
            status, 
            starts_at, 
            expires_at
        )
        VALUES (
            COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), 
            $2, 
            $3, 
            'active', 
            NOW(), 
            NOW() + ($4 * INTERVAL '1 day')
        )
        ON CONFLICT (user_id) DO UPDATE 
        SET 
            plan = EXCLUDED.plan,
            status = 'active',
            starts_at = CASE 
                WHEN subscriptions.expires_at <= NOW() THEN NOW() 
                ELSE subscriptions.starts_at 
            END,
            expires_at = GREATEST(subscriptions.expires_at, NOW()) + ($4 * INTERVAL '1 day'),
            updated_at = NOW()
        RETURNING subscription_id, user_id, plan, status, starts_at, expires_at, created_at, updated_at;
    `
	var output models.SubscriptionOutput

	err := r.DB.QueryRowContext(
		ctx,
		query,
		params.Id,
		params.UserId,
		params.Plan,
		params.Days,
	).Scan(
		&output.Id,
		&output.UserId,
		&output.Plan,
		&output.Status,
		&output.StartAt,
		&output.ExpiresAt,
		&output.CreatedAt,
		&output.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error upsert subscription: %w", err)
	}

	return &output, nil
}

func (r *subscriptionRepo) CheckPremiun(ctx context.Context, params *models.SubscriptionIsPremiunInput) (*models.SubscriptionIsPremiunOutput, error) {
	const query = `
    SELECT expires_at FROM subscriptions WHERE user_id = $1 AND status IN ('active', 'canceled') AND expires_at > NOW();`

	var output models.SubscriptionIsPremiunOutput

	err := r.DB.QueryRowContext(ctx, query, params.UserId).Scan(&output.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		output.Active = false
		output.ExpiresAt = time.Time{}
		return &output, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error verify status premiun: %w", err)
	}

	output.Active = true

	return &output, nil
}

func (r *subscriptionRepo) UpdateSubscriptionsStatusExpired(ctx context.Context, param *models.SubscriptionExpiredInput) (*models.SubscriptionExpiredOutput, error) {
	const query = `
        WITH target_rows AS (
            SELECT subscription_id
            FROM subscriptions
            WHERE status = 'active' AND expires_at <= NOW()
            LIMIT $1
            FOR UPDATE SKIP LOCKED
        )
        UPDATE subscriptions
        SET status = 'expired', updated_at = NOW()
        FROM target_rows
        WHERE subscriptions.subscription_id = target_rows.subscription_id;
    `
	result, err := r.DB.ExecContext(ctx, query, param.BatchSize)
	if err != nil {
		return nil, fmt.Errorf("r.DB.Exec: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("result.RowsAffected: %w", err) // Corrigido erro de digitação
	}

	out := models.SubscriptionExpiredOutput{
		RowsAffected: affected,
	}

	return &out, nil
}

func (r *subscriptionRepo) UpdateSubscriptionsStatusCancel(ctx context.Context, param *models.SubscriptionCancelInput) (*models.SubscriptionCancelOutput, error) {
	const query = `UPDATE subscriptions SET status = 'canceled',
    updated_at = NOW() WHERE user_id = $1 AND status = 'active'
    RETURNING subscription_id, user_id, status, updated_at;`

	var out models.SubscriptionCancelOutput
	err := r.DB.QueryRowContext(ctx, query, param.UserId).Scan(
		&out.SubscriptionId,
		&out.UserId,
		&out.Status,
		&out.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("r.DB.QueryRowContext: %w", err)
	}

	return &out, nil
}