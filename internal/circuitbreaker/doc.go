// Package circuitbreaker provides a thread-safe three-state circuit breaker
// for protecting downstream sinks in logpipe.
//
// States:
//
//	Closed   – normal operation; all calls are allowed through.
//	Open     – the downstream is considered unhealthy; calls are rejected
//	           immediately with ErrOpen until RecoveryTimeout elapses.
//	HalfOpen – a single probe call is allowed; success closes the circuit,
//	           failure re-opens it.
//
// Usage:
//
//	b := circuitbreaker.New(circuitbreaker.DefaultConfig())
//	if err := b.Allow(); err != nil {
//	    // fast-fail
//	}
//	// ... attempt operation ...
//	b.RecordSuccess() // or b.RecordFailure()
package circuitbreaker
