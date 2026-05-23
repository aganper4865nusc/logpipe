// Package aggregate implements time-window and count-based log line aggregation
// for logpipe.
//
// An Aggregator collects incoming log lines in an in-memory buffer and flushes
// them to a caller-supplied function when either:
//
//   - the buffer reaches maxSize lines, or
//   - a configurable interval elapses.
//
// Both conditions can be active simultaneously; whichever triggers first causes
// the flush.  A final flush is always performed when the context passed to Run
// is cancelled, ensuring no lines are lost on shutdown.
//
// Typical usage:
//
//	agg, err := aggregate.New(100, 5*time.Second, func(lines []string) {
//		// forward or process the batch
//	})
//	go agg.Run(ctx)
//	agg.Add(line)
package aggregate
