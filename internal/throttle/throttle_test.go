package throttle_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/throttle"
)

func TestNew_ZeroRate(t *testing.T) {
	_, _, err := throttle.New(func(string) error { return nil }, 0)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
}

func TestNew_NegativeRate(t *testing.T) {
	_, _, err := throttle.New(func(string) error { return nil }, -5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestWrite_ForwardsLine(t *testing.T) {
	var received string
	th, cancel, err := throttle.New(func(line string) error {
		received = line
		return nil
	}, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()

	if err := th.Write(context.Background(), "hello"); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if received != "hello" {
		t.Fatalf("expected 'hello', got %q", received)
	}
}

func TestWrite_ContextCancelled(t *testing.T) {
	// Use rate=1 so the token bucket drains after the first write.
	th, cancel, err := throttle.New(func(string) error { return nil }, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()

	// Consume the pre-filled token.
	_ = th.Write(context.Background(), "first")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer ctxCancel()

	err = th.Write(ctx, "second")
	if err == nil {
		t.Fatal("expected context deadline exceeded error")
	}
}

func TestWrite_ThrottlesCalls(t *testing.T) {
	const rate = 50
	var count atomic.Int64

	th, cancel, err := throttle.New(func(string) error {
		count.Add(1)
		return nil
	}, rate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()

	start := time.Now()
	for i := 0; i < 10; i++ {
		_ = th.Write(context.Background(), "line")
	}
	elapsed := time.Since(start)

	if count.Load() != 10 {
		t.Fatalf("expected 10 writes, got %d", count.Load())
	}
	// 10 lines at 50/s should take at least ~180 ms.
	if elapsed < 150*time.Millisecond {
		t.Fatalf("throttle too fast: elapsed %v", elapsed)
	}
}

func TestWrite_ThrottleCancelStopsRefill(t *testing.T) {
	th, cancel, err := throttle.New(func(string) error { return nil }, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cancel() // stop the refill goroutine immediately

	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer ctxCancel()

	// Drain the pre-filled token.
	_ = th.Write(ctx, "ok")

	err = th.Write(ctx, "blocked")
	if err == nil {
		t.Fatal("expected an error after throttle cancelled")
	}
}
