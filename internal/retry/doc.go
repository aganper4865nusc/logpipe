// Package retry implements a generic retry helper with exponential backoff
// and context-aware cancellation.
//
// Usage:
//
//	err := retry.Do(ctx, retry.DefaultConfig(), func() error {
//		return sink.Write(line)
//	})
//
// The helper backs off exponentially between attempts, capped at MaxDelay.
// If the context is cancelled at any point the operation returns immediately
// with the context error.
package retry
