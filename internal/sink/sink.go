package sink

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Sink represents a destination for log lines.
type Sink interface {
	Write(line string) error
	Close() error
}

// StdoutSink writes log lines to stdout.
type StdoutSink struct {
	w io.Writer
}

// NewStdoutSink creates a new StdoutSink.
func NewStdoutSink() *StdoutSink {
	return &StdoutSink{w: os.Stdout}
}

func (s *StdoutSink) Write(line string) error {
	_, err := fmt.Fprintln(s.w, line)
	return err
}

func (s *StdoutSink) Close() error { return nil }

// FileSink writes log lines to a file.
type FileSink struct {
	f *os.File
}

// NewFileSink opens or creates a file sink at the given path.
func NewFileSink(path string) (*FileSink, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("sink: open file %q: %w", path, err)
	}
	return &FileSink{f: f}, nil
}

func (s *FileSink) Write(line string) error {
	_, err := fmt.Fprintln(s.f, line)
	return err
}

func (s *FileSink) Close() error { return s.f.Close() }

// New creates a Sink based on the type string and target.
// Supported types: "stdout", "file".
func New(sinkType, target string) (Sink, error) {
	switch strings.ToLower(sinkType) {
	case "stdout":
		return NewStdoutSink(), nil
	case "file":
		if target == "" {
			return nil, fmt.Errorf("sink: file sink requires a non-empty target path")
		}
		return NewFileSink(target)
	default:
		return nil, fmt.Errorf("sink: unknown sink type %q", sinkType)
	}
}
