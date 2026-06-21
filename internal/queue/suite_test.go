package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nixmaldonado/blazeMailer/internal/model"
)

func sampleJob(id string) model.Job {
	return model.Job{ID: id, Email: model.EmailRequest{From: "a@b.com", To: "c@d.com"}, EnqueuedAt: time.Unix(0, 0)}
}

// RunQueueSuite asserts any Queue implementation honors the contract.
func RunQueueSuite(t *testing.T, newQ func() Queue) {
	t.Run("enqueue then dequeue FIFO", func(t *testing.T) {
		q := newQ()
		defer q.Close()
		ctx := context.Background()
		if err := q.Enqueue(ctx, sampleJob("1")); err != nil {
			t.Fatal(err)
		}
		if err := q.Enqueue(ctx, sampleJob("2")); err != nil {
			t.Fatal(err)
		}
		first, err := q.Dequeue(ctx)
		if err != nil || first.ID != "1" {
			t.Fatalf("got %v, %v; want job 1", first.ID, err)
		}
		second, err := q.Dequeue(ctx)
		if err != nil || second.ID != "2" {
			t.Fatalf("got %v, %v; want job 2", second.ID, err)
		}
	})

	t.Run("dequeue blocks then returns on ctx cancel", func(t *testing.T) {
		q := newQ()
		defer q.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := q.Dequeue(ctx)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("err = %v, want DeadlineExceeded", err)
		}
	})

	t.Run("depth reflects pending jobs", func(t *testing.T) {
		q := newQ()
		defer q.Close()
		ctx := context.Background()
		q.Enqueue(ctx, sampleJob("1"))
		n, err := q.Depth(ctx)
		if err != nil || n != 1 {
			t.Fatalf("depth = %d, %v; want 1", n, err)
		}
	})

	t.Run("dequeue after close returns ErrQueueClosed", func(t *testing.T) {
		q := newQ()
		ctx := context.Background()
		q.Close()
		if _, err := q.Dequeue(ctx); !errors.Is(err, ErrQueueClosed) {
			t.Fatalf("err = %v, want ErrQueueClosed", err)
		}
	})
}
