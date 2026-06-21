package queue

import (
	"context"
	"testing"
)

func TestMemoryQueue(t *testing.T) {
	RunQueueSuite(t, func() Queue { return NewMemory(100) })
}

func TestMemoryQueue_Full(t *testing.T) {
	q := NewMemory(1)
	defer q.Close()
	ctx := context.Background()
	if err := q.Enqueue(ctx, sampleJob("1")); err != nil {
		t.Fatal(err)
	}
	if err := q.Enqueue(ctx, sampleJob("2")); err != ErrQueueFull {
		t.Fatalf("err = %v, want ErrQueueFull", err)
	}
}
