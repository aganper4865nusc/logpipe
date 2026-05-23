// Package envelope defines the Envelope type — the fundamental unit of data
// that travels through the logpipe processing pipeline.
//
// An Envelope is created by a tailer when a new line is read from a source.
// It carries the raw line content alongside metadata such as the source name,
// read timestamp, byte offset, and an extensible field map that transforms
// may populate (e.g. AddField, AddTimestamp).
//
// Envelope values are treated as immutable once created; all mutation helpers
// (WithField, Clone) return new copies to avoid data races in concurrent
// pipeline stages.
package envelope
