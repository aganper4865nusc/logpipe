package sink

import (
	"context"
	"fmt"

	"github.com/yourorg/logpipe/internal/retry"
)

// RetrySink wraps another Sink and retries failed writes according to a
// retry.Config. It is transparent to the caller: a successful write within
// the allowed attempts returns nil.
type RetrySink struct {
	inner  Sink
	cfg    retry.Config
	ctx    context.Context
}

// NewRetrySink wraps inner with retry logic using cfg.
// ctx controls the lifetime of retry back-offs.
func NewRetrySink(ctx context.Context, inner Sink, cfg retry.Config) *RetrySink {
	return &RetrySink{inner: inner, cfg: cfg, ctx: ctx}
}

// Write attempts to write line to the inner sink, retrying on failure.
func (r *RetrySink) Write(line string) error {
	err := retry.Do(r.ctx, r.cfg, func() error {
		return r.inner.Write(line)
	})
	if err != nil {
		return fmt.Errorf("retry_sink: %w", err)
	}
	return nil
}

// Close closes the inner sink.
func (r *RetrySink) Close() error {
	return r.inner.Close()
}
