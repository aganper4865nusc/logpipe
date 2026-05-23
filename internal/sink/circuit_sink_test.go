package sink_test

import (
	"errors"
	"testing"
	"time"

	"github.com/logpipe/logpipe/internal/circuitbreaker"
	"github.com/logpipe/logpipe/internal/sink"
)

type failingSink struct{ calls int; failUntil int }

func (f *failingSink) Write(_ string) error {
	f.calls++
	if f.calls <= f.failUntil {
		return errors.New("downstream error")
	}
	return nil
}

func fastCBConfig(threshold int) circuitbreaker.Config {
	return circuitbreaker.Config{
		FailureThreshold: threshold,
		RecoveryTimeout:  10 * time.Millisecond,
	}
}

func TestCircuitSink_PassesOnSuccess(t *testing.T) {
	inner := &failingSink{failUntil: 0}
	cs := sink.NewCircuitSink(inner, fastCBConfig(3))
	if err := cs.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cs.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed, got %v", cs.State())
	}
}

func TestCircuitSink_OpensAfterThreshold(t *testing.T) {
	inner := &failingSink{failUntil: 3}
	cs := sink.NewCircuitSink(inner, fastCBConfig(3))
	for i := 0; i < 3; i++ {
		_ = cs.Write("line")
	}
	if cs.State() != circuitbreaker.StateOpen {
		t.Fatalf("expected Open, got %v", cs.State())
	}
	err := cs.Write("rejected")
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
	if !errors.Is(err, circuitbreaker.ErrOpen) {
		t.Fatalf("expected ErrOpen wrapped, got: %v", err)
	}
}

func TestCircuitSink_RecoveryAfterTimeout(t *testing.T) {
	inner := &failingSink{failUntil: 2}
	cs := sink.NewCircuitSink(inner, fastCBConfig(2))
	_ = cs.Write("fail1")
	_ = cs.Write("fail2")
	time.Sleep(20 * time.Millisecond)
	// probe should succeed now (inner no longer fails)
	if err := cs.Write("probe"); err != nil {
		t.Fatalf("expected recovery, got: %v", err)
	}
	if cs.State() != circuitbreaker.StateClosed {
		t.Fatalf("expected Closed after recovery, got %v", cs.State())
	}
}
