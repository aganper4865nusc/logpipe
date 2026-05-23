// Package envelope defines the core log record type that flows through the pipeline.
// Each line read from a source is wrapped in an Envelope before being forwarded
// to filters, transforms, and sinks.
package envelope

import "time"

// Envelope wraps a raw log line with metadata about its origin.
type Envelope struct {
	// Source is the name of the tailer or input that produced this line.
	Source string

	// Line is the raw log line content.
	Line string

	// Timestamp is the wall-clock time at which the line was read.
	Timestamp time.Time

	// Offset is the byte offset of the line within the source file, if applicable.
	Offset int64

	// Fields holds any additional key/value metadata attached by transforms.
	Fields map[string]string
}

// New creates an Envelope with the current UTC time.
func New(source, line string, offset int64) Envelope {
	return Envelope{
		Source:    source,
		Line:      line,
		Timestamp: time.Now().UTC(),
		Offset:    offset,
		Fields:    make(map[string]string),
	}
}

// WithField returns a shallow copy of the Envelope with the given field set.
func (e Envelope) WithField(key, value string) Envelope {
	copy := e
	copy.Fields = make(map[string]string, len(e.Fields)+1)
	for k, v := range e.Fields {
		copy.Fields[k] = v
	}
	copy.Fields[key] = value
	return copy
}

// Clone returns a deep copy of the Envelope.
func (e Envelope) Clone() Envelope {
	return e.WithField("", "") // reuses copy logic; empty key is harmless
}
