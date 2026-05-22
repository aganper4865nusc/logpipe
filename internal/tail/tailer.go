package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single tailed log line with its source.
type Line struct {
	Source string
	Text   string
	Time   time.Time
}

// Tailer tails a file and emits lines to an output channel.
type Tailer struct {
	path   string
	out    chan<- Line
	pollInterval time.Duration
}

// New creates a new Tailer for the given file path.
func New(path string, out chan<- Line) *Tailer {
	return &Tailer{
		path:         path,
		out:          out,
		pollInterval: 250 * time.Millisecond,
	}
}

// Run opens the file, seeks to the end, and streams new lines until ctx is cancelled.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// No new data yet; wait before retrying.
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(t.pollInterval):
				}
				continue
			}
			return err
		}

		if len(line) > 0 {
			t.out <- Line{
				Source: t.path,
				Text:   line,
				Time:   time.Now().UTC(),
			}
		}
	}
}
