// Package multiline provides a log line assembler that coalesces
// continuation lines (e.g. stack traces) into a single logical event.
package multiline

import (
	"regexp"
	"strings"
	"time"
)

// Config controls how continuation lines are detected.
type Config struct {
	// StartPattern marks the beginning of a new logical line.
	// Any line that does NOT match is treated as a continuation.
	StartPattern string
	// MaxLines is the maximum number of raw lines to fold into one event.
	// 0 means no limit.
	MaxLines int
	// FlushTimeout is the maximum time to wait for more continuation lines
	// before emitting the current buffer. 0 disables the timeout.
	FlushTimeout time.Duration
}

// Assembler folds multi-line log events into single strings.
type Assembler struct {
	start   *regexp.Regexp
	maxLine int
	timeout time.Duration
	buf     []string
}

// New creates an Assembler from cfg.
// Returns an error if StartPattern is not a valid regular expression.
func New(cfg Config) (*Assembler, error) {
	re, err := regexp.Compile(cfg.StartPattern)
	if err != nil {
		return nil, err
	}
	return &Assembler{
		start:   re,
		maxLine: cfg.MaxLines,
		timeout: cfg.FlushTimeout,
	}, nil
}

// Feed accepts a raw line and returns zero or one completed logical events.
// A completed event is returned when a new start-pattern is detected or
// the internal buffer reaches MaxLines.
func (a *Assembler) Feed(line string) (event string, ready bool) {
	isStart := a.start.MatchString(line)

	if isStart && len(a.buf) > 0 {
		event = a.flush()
		a.buf = append(a.buf, line)
		return event, true
	}

	a.buf = append(a.buf, line)

	if a.maxLine > 0 && len(a.buf) >= a.maxLine {
		return a.flush(), true
	}

	return "", false
}

// Flush forces the current buffer to be emitted regardless of state.
// Returns an empty string if the buffer is empty.
func (a *Assembler) Flush() string {
	return a.flush()
}

func (a *Assembler) flush() string {
	if len(a.buf) == 0 {
		return ""
	}
	out := strings.Join(a.buf, "\n")
	a.buf = a.buf[:0]
	return out
}
