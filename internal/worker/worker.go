package worker

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/nixmaldonado/blazeMailer/internal/delivery"
	"github.com/nixmaldonado/blazeMailer/internal/queue"
)

type Pool struct {
	q    queue.Queue
	s    delivery.Sender
	size int
}

func NewPool(q queue.Queue, s delivery.Sender, size int) *Pool {
	return &Pool{q: q, s: s, size: size}
}

func (p *Pool) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < p.size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.loop(ctx)
		}()
	}
	wg.Wait()
}

func (p *Pool) loop(ctx context.Context) {
	for {
		job, err := p.q.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) ||
				errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, queue.ErrQueueClosed) {
				return
			}
			continue
		}
		if err := p.s.Send(ctx, job.Email); err != nil {
			log.Printf("delivery failed job=%s: %v", job.ID, err)
		}
	}
}
