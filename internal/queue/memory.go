package queue

import (
	"context"
	"log"
	"sync"

	"github.com/nixmaldonado/blazeMailer/internal/model"
)

type Memory struct {
	jobs   chan model.Job
	mu     sync.Mutex
	closed bool
}

func NewMemory(capacity int) *Memory {
	return &Memory{jobs: make(chan model.Job, capacity)}
}

func (m *Memory) Enqueue(_ context.Context, j model.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return ErrQueueClosed
	}
	select {
	case m.jobs <- j:
		return nil
	default:
		return ErrQueueFull
	}
}

func (m *Memory) Dequeue(ctx context.Context) (model.Job, error) {
	select {
	case j, ok := <-m.jobs:
		if !ok {
			return model.Job{}, ErrQueueClosed
		}
		return j, nil
	case <-ctx.Done():
		return model.Job{}, ctx.Err()
	}
}

func (m *Memory) DeadLetter(_ context.Context, j model.Job, reason string) error {
	log.Printf("DLQ job=%s reason=%s", j.ID, reason)
	return nil
}

func (m *Memory) Depth(_ context.Context) (int, error) { return len(m.jobs), nil }

func (m *Memory) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		m.closed = true
		close(m.jobs)
	}
	return nil
}
