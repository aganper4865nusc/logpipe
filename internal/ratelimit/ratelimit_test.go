package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/ratelimit"
)

func TestLimiter_TryConsume_BurstAvailable(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Rate: 1, Burst: 5})
	for i := 0; i < 5; i++ {
		if !l.TryConsume() {
			t.Fatalf("expected token %d to be available", i)
		}
	}
	if l.TryConsume() {
		t.Fatal("expected no token after burst exhausted")
	}
}

func TestLimiter_TryConsume_RefillsOverTime(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Rate: 1000, Burst: 1})
	if !l.TryConsume() {
		t.Fatal("expected first token")
	}
	// At 1000 tokens/sec, 10ms should refill ~10 tokens
	time.Sleep(10 * time.Millisecond)
	if !l.TryConsume() {
		t.Fatal("expected token after sleep")
	}
}

func TestLimiter_Wait_ConsumesToken(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Rate: 10000, Burst: 10})
	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLimiter_Wait_ContextCancelled(t *testing.T) {
	// Zero rate: tokens never refill after burst is consumed
	l := ratelimit.New(ratelimit.Config{Rate: 0, Burst: 0})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := l.Wait(ctx)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}

func TestLimiter_Wait_ContextAlreadyCancelled(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{Rate: 1, Burst: 1})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := l.Wait(ctx); err == nil {
		t.Fatal("expected error for already-cancelled context")
	}
}
