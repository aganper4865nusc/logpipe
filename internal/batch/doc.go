// Package batch implements size- and time-based log line batching for logpipe.
//
// A Batcher collects individual log lines sent via Send and periodically
// flushes them to a Flusher callback. Flushing is triggered by whichever
// condition occurs first:
//
//   - The internal buffer reaches the configured maximum size (maxSize).
//   - The configured timeout elapses since the last flush.
//   - The context passed to Run is cancelled (remaining lines are flushed
//     before the goroutine exits).
//
// Typical usage:
//
//	b := batch.New(50, 500*time.Millisecond, func(lines []string) {
//	    sink.WriteAll(lines)
//	})
//	go b.Run(ctx)
//	for line := range source {
//	    b.Send(line)
//	}
package batch
