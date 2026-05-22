// Package buffer implements a bounded ring buffer for decoupling log
// producers from consumers in the logpipe pipeline.
//
// # Overview
//
// A RingBuffer has a fixed capacity set at construction time. When the buffer
// is full two behaviours are available:
//
//   - dropOnFull=false: Push returns ErrFull and the caller decides what to do.
//   - dropOnFull=true:  The oldest entry is silently evicted and a dropped
//     counter is incremented, allowing the pipeline to keep moving
//     under back-pressure.
//
// All operations are safe for concurrent use.
package buffer
