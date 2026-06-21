package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nixmaldonado/blazeMailer/internal/model"
	"github.com/nixmaldonado/blazeMailer/internal/queue"
)

type countingSender struct {
	mu    sync.Mutex
	count int
}

func (c *countingSender) Send(_ context.Context, _ model.EmailRequest) error {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
	return nil
}

func TestPool_DeliversAllJobs(t *testing.T) {
	q := queue.NewMemory(100)
	cs := &countingSender{}
	for i := 0; i < 10; i++ {
		q.Enqueue(context.Background(), model.Job{ID: string(rune('a' + i))})
	}
	ctx, cancel := context.WithCancel(context.Background())
	p := NewPool(q, cs, 4)
	done := make(chan struct{})
	go func() { p.Run(ctx); close(done) }()

	deadline := time.After(2 * time.Second)
	for {
		cs.mu.Lock()
		n := cs.count
		cs.mu.Unlock()
		if n == 10 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("only %d delivered", n)
		case <-time.After(5 * time.Millisecond):
		}
	}
	cancel()
	<-done
}

func TestPool_DrainsOnQueueClose(t *testing.T) {
	q := queue.NewMemory(100)
	cs := &countingSender{}
	for i := 0; i < 10; i++ {
		q.Enqueue(context.Background(), model.Job{ID: string(rune('a' + i))})
	}
	q.Close() // close after enqueue: workers must drain all 10, then exit

	p := NewPool(q, cs, 4)
	done := make(chan struct{})
	go func() { p.Run(context.Background()); close(done) }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("pool did not drain and exit after queue close")
	}
	cs.mu.Lock()
	got := cs.count
	cs.mu.Unlock()
	if got != 10 {
		t.Fatalf("delivered %d jobs, want 10", got)
	}
}
