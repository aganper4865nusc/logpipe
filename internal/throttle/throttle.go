// Package throttle provides a sink wrapper that enforces a maximum
// number of lines forwarded per second across all configured sinks.
package throttle

import (
	"context"
	"errors"
	"time"
)

// ErrZeroRate is returned when a rate of zero or less is configured.
var ErrZeroRate = errors.New("throttle: rate must be greater than zero")

// Sink is a function type that accepts a line string.
type Sink func(line string) error

// Throttle wraps a Sink and limits throughput to at most Rate lines per second.
type Throttle struct {
	sink     Sink
	rate     int
	ticker   *time.Ticker
	tokens   chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a Throttle that allows at most ratePerSec lines per second
// to pass through to sink. The background token-refill goroutine runs until
// the returned cancel function is called.
func New(sink Sink, ratePerSec int) (*Throttle, context.CancelFunc, error) {
	if ratePerSec <= 0 {
		return nil, nil, ErrZeroRate
	}

	ctx, cancel := context.WithCancel(context.Background())

	t := &Throttle{
		sink:   sink,
		rate:   ratePerSec,
		ticker: time.NewTicker(time.Second / time.Duration(ratePerSec)),
		tokens: make(chan struct{}, ratePerSec),
		ctx:    ctx,
		cancel: cancel,
	}

	// Pre-fill one token so the first call is immediate.
	t.tokens <- struct{}{}

	go t.refill()

	return t, cancel, nil
}

// refill pumps one token into the channel on every tick.
func (t *Throttle) refill() {
	for {
		select {
		case <-t.ctx.Done():
			t.ticker.Stop()
			return
		case <-t.ticker.C:
			select {
			case t.tokens <- struct{}{}:
			default: // bucket full, discard
			}
		}
	}
}

// Write blocks until a token is available, then forwards line to the
// underlying sink. Returns an error if the context is cancelled while waiting.
func (t *Throttle) Write(ctx context.Context, line string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.ctx.Done():
		return t.ctx.Err()
	case <-t.tokens:
		return t.sink(line)
	}
}
