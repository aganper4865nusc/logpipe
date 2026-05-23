package backpressure_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/backpressure"
)

func TestNew_InvalidCapacity(t *testing.T) {
	_, err := backpressure.New(0)
	if err == nil {
		t.Fatal("expected error for zero capacity")
	}
}

func TestSend_DeliversLine(t *testing.T) {
	ctrl, err := backpressure.New(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer ctrl.Close()

	ctx := context.Background()
	if err := ctrl.Send(ctx, "hello"); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	select {
	case got := <-ctrl.Receive():
		if got != "hello" {
			t.Fatalf("want %q, got %q", "hello", got)
		}
	default:
		t.Fatal("expected a message in the channel")
	}
}

func TestTrySend_DropsWhenFull(t *testing.T) {
	ctrl, _ := backpressure.New(2)
	defer ctrl.Close()

	_ = ctrl.TrySend("a")
	_ = ctrl.TrySend("b")

	err := ctrl.TrySend("c")
	if err != backpressure.ErrDropped {
		t.Fatalf("want ErrDropped, got %v", err)
	}
	if ctrl.Dropped() != 1 {
		t.Fatalf("want 1 dropped, got %d", ctrl.Dropped())
	}
}

func TestSend_RespectsContextCancellation(t *testing.T) {
	ctrl, _ := backpressure.New(1)
	defer ctrl.Close()

	// fill the buffer
	_ = ctrl.TrySend("fill")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := ctrl.Send(ctx, "blocked")
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestDropped_CountsAccumulate(t *testing.T) {
	ctrl, _ := backpressure.New(1)
	defer ctrl.Close()

	_ = ctrl.TrySend("keep")
	for i := 0; i < 5; i++ {
		_ = ctrl.TrySend("drop")
	}
	if ctrl.Dropped() != 5 {
		t.Fatalf("want 5 dropped, got %d", ctrl.Dropped())
	}
}
