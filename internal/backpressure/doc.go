// Package backpressure implements token-based flow control for logpipe
// pipelines.
//
// A Controller wraps a buffered channel and exposes two sending modes:
//
//   - Send: blocks the caller until capacity is available or the context is
//     cancelled. Suitable for sources that must not lose data.
//
//   - TrySend: returns immediately with ErrDropped if the buffer is full.
//     Suitable for high-throughput sources where occasional loss is acceptable.
//
// Dropped message counts are tracked atomically and can be scraped by the
// metrics registry for alerting.
package backpressure
