package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/circuitbreaker"
)

func TestBreaker_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(circuitbreaker.DefaultConfig())
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestBreaker_OpensAfterThreshold(t *testing.T) {
	cfg := circuitbreaker.Config{FailureThreshold: 3, RecoveryTimeout: time.Minute}
	b := circuitbreaker.New(cfg)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected Open, got %v", b.State())
	}
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_SuccessResetsClosed(t *testing.T) {
	cfg := circuitbreaker.Config{FailureThreshold: 2, RecoveryTimeout: time.Minute}
	b := circuitbreaker.New(cfg)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after success, got %v", b.State())
	}
}

func TestBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cfg := circuitbreaker.Config{FailureThreshold: 1, RecoveryTimeout: 10 * time.Millisecond}
	b := circuitbreaker.New(cfg)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in HalfOpen probe, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", b.State())
	}
}

func TestBreaker_HalfOpenSuccessCloses(t *testing.T) {
	cfg := circuitbreaker.Config{FailureThreshold: 1, RecoveryTimeout: 5 * time.Millisecond}
	b := circuitbreaker.New(cfg)
	b.RecordFailure()
	time.Sleep(10 * time.Millisecond)
	_ = b.Allow() // transition to HalfOpen
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after HalfOpen success, got %v", b.State())
	}
}

func TestBreaker_DefaultConfigSanity(t *testing.T) {
	cfg := circuitbreaker.DefaultConfig()
	if cfg.FailureThreshold <= 0 {
		t.Fatal("FailureThreshold must be positive")
	}
	if cfg.RecoveryTimeout <= 0 {
		t.Fatal("RecoveryTimeout must be positive")
	}
}
