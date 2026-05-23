// Package batch provides a size- and time-based line batcher for logpipe.
// Lines are accumulated and flushed either when the batch reaches its
// maximum size or when a configurable timeout elapses, whichever comes first.
package batch

import (
	"context"
	"time"
)

// Flusher is called with a completed batch of lines.
type Flusher func(lines []string)

// Batcher accumulates lines and flushes them in groups.
type Batcher struct {
	maxSize  int
	timeout  time.Duration
	flusher  Flusher
	input    chan string
}

// New creates a Batcher that flushes when maxSize lines are buffered or
// timeout elapses. Call Run to start processing.
func New(maxSize int, timeout time.Duration, flusher Flusher) *Batcher {
	if maxSize <= 0 {
		maxSize = 100
	}
	if timeout <= 0 {
		timeout = time.Second
	}
	return &Batcher{
		maxSize: maxSize,
		timeout: timeout,
		flusher: flusher,
		input:   make(chan string, maxSize*2),
	}
}

// Send enqueues a line for batching. It is non-blocking; lines are dropped
// if the internal channel is full.
func (b *Batcher) Send(line string) {
	select {
	case b.input <- line:
	default:
	}
}

// Run starts the batching loop and blocks until ctx is cancelled.
func (b *Batcher) Run(ctx context.Context) {
	buf := make([]string, 0, b.maxSize)
	ticker := time.NewTicker(b.timeout)
	defer ticker.Stop()

	flush := func() {
		if len(buf) == 0 {
			return
		}
		copy := make([]string, len(buf))
		copy(copy, buf)
		b.flusher(copy)
		buf = buf[:0]
	}

	for {
		select {
		case line := <-b.input:
			buf = append(buf, line)
			if len(buf) >= b.maxSize {
				flush()
				ticker.Reset(b.timeout)
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush()
			return
		}
	}
}
