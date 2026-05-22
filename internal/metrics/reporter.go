package metrics

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Reporter periodically writes a metrics summary to a writer.
type Reporter struct {
	reg      *Registry
	out      io.Writer
	interval time.Duration
}

// NewReporter creates a Reporter that writes to out every interval.
func NewReporter(reg *Registry, out io.Writer, interval time.Duration) *Reporter {
	return &Reporter{
		reg:      reg,
		out:      out,
		interval: interval,
	}
}

// Run starts the periodic reporting loop; it returns when ctx is cancelled.
func (r *Reporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			r.report()
			return
		case <-ticker.C:
			r.report()
		}
	}
}

func (r *Reporter) report() {
	snap := r.reg.Snapshot()
	fmt.Fprintf(
		r.out,
		"[metrics] read=%d filtered=%d forwarded=%d sink_errors=%d\n",
		snap.LinesRead,
		snap.LinesFiltered,
		snap.LinesForwarded,
		snap.SinkErrors,
	)
}
