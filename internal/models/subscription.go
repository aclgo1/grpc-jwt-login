package models

import "time"

type SubscriptionInput struct {
	Id string
	UserId string
	Plan string
	Days int64
}

type SubscriptionOutput struct {
	Id string `db:"subscription_id"`
	UserId string  `db:"user_id"`
	Plan string `db:"plan"`
	Status string `db:"status"`
	StartAt time.Time `db:"starts_at"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SubscriptionIsPremiunInput struct{
	UserId string
}

type SubscriptionIsPremiunOutput struct{
	Active bool
	ExpiresAt time.Time
}

type SubscriptionExpiredInput struct {
	BatchSize int64
}

type SubscriptionExpiredOutput struct {
	RowsAffected int64
}

type SubscriptionCancelInput struct {
	UserId string
}

type SubscriptionCancelOutput struct {
	SubscriptionId string    `db:"subscription_id"`
	UserId         string    `db:"user_id"`
	Status         string    `db:"status"`
	UpdatedAt      time.Time `db:"updated_at"`
}