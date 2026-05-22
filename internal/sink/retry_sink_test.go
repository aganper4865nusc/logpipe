package sink_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/retry"
	"github.com/yourorg/logpipe/internal/sink"
)

type countingSink struct {
	calls  int32
	failN  int32 // fail this many times before succeeding
	closed bool
}

func (c *countingSink) Write(_ string) error {
	n := atomic.AddInt32(&c.calls, 1)
	if n <= atomic.LoadInt32(&c.failN) {
		return errors.New("sink unavailable")
	}
	return nil
}

func (c *countingSink) Close() error {
	c.closed = true
	return nil
}

func fastRetry(max int) retry.Config {
	return retry.Config{
		MaxAttempts:  max,
		InitialDelay: time.Millisecond,
		MaxDelay:     5 * time.Millisecond,
		Multiplier:   2.0,
	}
}

func TestRetrySink_SucceedsAfterRetries(t *testing.T) {
	inner := &countingSink{failN: 2}
	s := sink.NewRetrySink(context.Background(), inner, fastRetry(5))
	if err := s.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetrySink_FailsAfterMaxAttempts(t *testing.T) {
	inner := &countingSink{failN: 10}
	s := sink.NewRetrySink(context.Background(), inner, fastRetry(3))
	err := s.Write("hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRetrySink_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	inner := &countingSink{failN: 5}
	s := sink.NewRetrySink(ctx, inner, fastRetry(5))
	err := s.Write("hello")
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

func TestRetrySink_Close(t *testing.T) {
	inner := &countingSink{}
	s := sink.NewRetrySink(context.Background(), inner, fastRetry(3))
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !inner.closed {
		t.Fatal("expected inner sink to be closed")
	}
}
