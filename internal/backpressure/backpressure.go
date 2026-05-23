// Package backpressure provides a simple token-based flow-control mechanism
// that allows producers to block or drop messages when a downstream consumer
// is unable to keep up.
package backpressure

import (
	"context"
	"errors"
	"sync/atomic"
)

// ErrDropped is returned by TrySend when the internal channel is full and
// the controller is configured to drop rather than block.
var ErrDropped = errors.New("backpressure: message dropped")

// Controller manages flow between a producer and a consumer.
type Controller struct {
	ch      chan string
	dropped atomic.Int64
}

// New creates a Controller with the given buffer capacity.
// capacity must be greater than zero.
func New(capacity int) (*Controller, error) {
	if capacity <= 0 {
		return nil, errors.New("backpressure: capacity must be > 0")
	}
	return &Controller{ch: make(chan string, capacity)}, nil
}

// Send blocks until the line can be placed in the buffer or ctx is cancelled.
func (c *Controller) Send(ctx context.Context, line string) error {
	select {
	case c.ch <- line:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TrySend attempts a non-blocking send. If the buffer is full ErrDropped is
// returned and the internal drop counter is incremented.
func (c *Controller) TrySend(line string) error {
	select {
	case c.ch <- line:
		return nil
	default:
		c.dropped.Add(1)
		return ErrDropped
	}
}

// Receive returns the read-only channel that consumers should range over.
func (c *Controller) Receive() <-chan string {
	return c.ch
}

// Dropped returns the total number of lines dropped since creation.
func (c *Controller) Dropped() int64 {
	return c.dropped.Load()
}

// Close signals that no more lines will be sent.
func (c *Controller) Close() {
	close(c.ch)
}
