package model

import "time"

type Job struct {
	ID             string       `json:"id"`
	IdempotencyKey string       `json:"idempotency_key,omitempty"`
	Email          EmailRequest `json:"email"`
	Attempts       int          `json:"attempts"`
	EnqueuedAt     time.Time    `json:"enqueued_at"`
}
