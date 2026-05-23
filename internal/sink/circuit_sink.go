package sink

import (
	"fmt"

	"github.com/logpipe/logpipe/internal/circuitbreaker"
)

// CircuitSink wraps a Sink with a circuit breaker so that repeated downstream
// failures stop propagating until the circuit recovers.
type CircuitSink struct {
	inner   Sink
	breaker *circuitbreaker.Breaker
}

// NewCircuitSink creates a CircuitSink wrapping inner with the given breaker
// configuration.
func NewCircuitSink(inner Sink, cfg circuitbreaker.Config) *CircuitSink {
	return &CircuitSink{
		inner:   inner,
		breaker: circuitbreaker.New(cfg),
	}
}

// Write checks the circuit before forwarding the line to the inner sink.
// On success it records a success; on failure it records a failure and returns
// the error.
func (c *CircuitSink) Write(line string) error {
	if err := c.breaker.Allow(); err != nil {
		return fmt.Errorf("circuit open: %w", err)
	}
	if err := c.inner.Write(line); err != nil {
		c.breaker.RecordFailure()
		return err
	}
	c.breaker.RecordSuccess()
	return nil
}

// State exposes the underlying breaker state for observability.
func (c *CircuitSink) State() circuitbreaker.State {
	return c.breaker.State()
}
