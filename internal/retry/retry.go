// Package retry provides configurable retry logic with exponential backoff
// for use in sink writers and other I/O-bound operations.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Config holds retry policy parameters.
type Config struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Do executes fn up to cfg.MaxAttempts times, backing off between attempts.
// It returns nil on the first success, or ErrMaxAttemptsReached if all fail.
// The context is checked before each attempt.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	delay := cfg.InitialDelay
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := fn(); err == nil {
			return nil
		}
		if attempt == cfg.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		next := time.Duration(math.Min(
			float64(delay)*cfg.Multiplier,
			float64(cfg.MaxDelay),
		))
		delay = next
	}
	return ErrMaxAttemptsReached
}
