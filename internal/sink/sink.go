package sink

import (
	"fmt"
	"os"
)

// StdoutSink writes log lines to standard output.
type StdoutSink struct{}

// NewStdoutSink returns a new StdoutSink.
func NewStdoutSink() *StdoutSink {
	return &StdoutSink{}
}

// Write prints line followed by a newline to stdout.
func (s *StdoutSink) Write(line string) error {
	_, err := fmt.Fprintln(os.Stdout, line)
	return err
}

// FileSink writes log lines to a file on disk.
type FileSink struct {
	f *os.File
}

// NewFileSink opens (or creates) the file at path for appending.
func NewFileSink(path string) (*FileSink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("sink: open file %q: %w", path, err)
	}
	return &FileSink{f: f}, nil
}

// Write appends line followed by a newline to the file.
func (s *FileSink) Write(line string) error {
	_, err := fmt.Fprintln(s.f, line)
	return err
}

// Close releases the underlying file handle.
func (s *FileSink) Close() error {
	return s.f.Close()
}

// New is a factory that returns a Sink based on the type string.
// Supported types: "stdout", "file" (requires target field).
func New(sinkType, target string) (interface{ Write(string) error }, error) {
	switch sinkType {
	case "stdout":
		return NewStdoutSink(), nil
	case "file":
		if target == "" {
			return nil, fmt.Errorf("sink: file sink requires a target path")
		}
		return NewFileSink(target)
	default:
		return nil, fmt.Errorf("sink: unknown sink type %q", sinkType)
	}
}
