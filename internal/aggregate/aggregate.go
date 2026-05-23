// Package aggregate provides time-window and count-based aggregation of log lines.
// It collects incoming lines and emits a summary when a flush condition is met.
package aggregate

import (
	"context"
	"sync"
	"time"
)

// Func is called with the accumulated lines when a flush occurs.
type Func func(lines []string)

// Aggregator buffers lines and flushes them on size or interval.
type Aggregator struct {
	mu       sync.Mutex
	buf      []string
	maxSize  int
	interval time.Duration
	flushFn  Func
}

// New creates an Aggregator. maxSize <= 0 disables size-based flushing.
// interval <= 0 disables time-based flushing.
func New(maxSize int, interval time.Duration, fn Func) (*Aggregator, error) {
	if fn == nil {
		return nil, errNilFunc
	}
	if maxSize <= 0 && interval <= 0 {
		return nil, errNoFlushCondition
	}
	return &Aggregator{
		maxSize:  maxSize,
		interval: interval,
		flushFn:  fn,
	}, nil
}

// Add appends a line to the buffer, flushing if maxSize is reached.
func (a *Aggregator) Add(line string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buf = append(a.buf, line)
	if a.maxSize > 0 && len(a.buf) >= a.maxSize {
		a.flush()
	}
}

// Run starts the interval-based flush loop. It blocks until ctx is cancelled.
func (a *Aggregator) Run(ctx context.Context) {
	if a.interval <= 0 {
		<-ctx.Done()
		a.Flush()
		return
	}
	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.Flush()
		case <-ctx.Done():
			a.Flush()
			return
		}
	}
}

// Flush immediately drains the buffer and calls the flush function.
func (a *Aggregator) Flush() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.flush()
}

// flush must be called with a.mu held.
func (a *Aggregator) flush() {
	if len(a.buf) == 0 {
		return
	}
	copy := make([]string, len(a.buf))
	for i, v := range a.buf {
		copy[i] = v
	}
	a.buf = a.buf[:0]
	a.flushFn(copy)
}
