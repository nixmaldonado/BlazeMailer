package queue

import (
	"context"
	"errors"

	"github.com/nixmaldonado/blazeMailer/internal/model"
)

var (
	ErrQueueFull   = errors.New("queue: full")
	ErrQueueClosed = errors.New("queue: closed")
)

type Queue interface {
	Enqueue(ctx context.Context, j model.Job) error
	Dequeue(ctx context.Context) (model.Job, error)
	DeadLetter(ctx context.Context, j model.Job, reason string) error
	Depth(ctx context.Context) (int, error)
	Close() error
}
