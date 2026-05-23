// Package circuitbreaker implements a simple three-state circuit breaker
// (Closed, Open, HalfOpen) to protect downstream sinks from repeated failures.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current circuit breaker state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing fast
	StateHalfOpen              // probing recovery
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuitbreaker: circuit is open")

// Config holds tuneable parameters for the breaker.
type Config struct {
	// FailureThreshold is the number of consecutive failures before opening.
	FailureThreshold int
	// RecoveryTimeout is how long the circuit stays open before probing.
	RecoveryTimeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: 5,
		RecoveryTimeout:  30 * time.Second,
	}
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	cfg        Config
	mu         sync.Mutex
	state      State
	failures   int
	lastOpened time.Time
}

// New creates a Breaker with the given config.
func New(cfg Config) *Breaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = DefaultConfig().FailureThreshold
	}
	if cfg.RecoveryTimeout <= 0 {
		cfg.RecoveryTimeout = DefaultConfig().RecoveryTimeout
	}
	return &Breaker{cfg: cfg}
}

// Allow returns nil if the call should proceed, or ErrOpen if it should be
// rejected. It transitions the breaker to HalfOpen when the recovery timeout
// has elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.lastOpened) >= b.cfg.RecoveryTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the breaker to Closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.cfg.FailureThreshold {
		b.state = StateOpen
		b.lastOpened = time.Now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
