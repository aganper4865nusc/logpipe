// Package ratelimit provides a token-bucket rate limiter for controlling
// log line throughput through the pipeline.
package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Limiter controls the rate at which log lines are processed.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// Config holds rate limiter configuration.
type Config struct {
	// Rate is the number of lines allowed per second.
	Rate float64
	// Burst is the maximum number of tokens that can accumulate.
	Burst float64
}

// New creates a Limiter with the given config.
func New(cfg Config) *Limiter {
	now := time.Now()
	return &Limiter{
		tokens:   cfg.Burst,
		max:      cfg.Burst,
		rate:     cfg.Rate,
		lastTick: now,
		clock:    time.Now,
	}
}

// Wait blocks until a token is available or ctx is cancelled.
// Returns ctx.Err() if the context is cancelled before a token is available.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if l.tryConsume() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 5):
		}
	}
}

// TryConsume attempts to consume a token without blocking.
// Returns true if a token was consumed.
func (l *Limiter) TryConsume() bool {
	return l.tryConsume()
}

func (l *Limiter) tryConsume() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now
	l.tokens = min(l.max, l.tokens+elapsed*l.rate)
	if l.tokens >= 1.0 {
		l.tokens--
		return true
	}
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
