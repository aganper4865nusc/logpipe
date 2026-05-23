// Package rotator detects when a tailed file has been rotated (truncated or
// replaced) and signals the caller so the tailer can re-open the source.
package rotator

import (
	"errors"
	"os"
	"time"
)

// ErrRotated is returned by Check when a rotation is detected.
var ErrRotated = errors.New("rotator: file has been rotated")

// Rotator watches a file for truncation or inode change.
type Rotator struct {
	path     string
	inode    uint64
	size     int64
	interval time.Duration
}

// Config holds options for New.
type Config struct {
	// Interval between rotation checks. Defaults to 5 seconds.
	Interval time.Duration
}

// New creates a Rotator for the given path, capturing the current inode and
// size as the baseline. Returns an error if the file cannot be stat-ed.
func New(path string, cfg Config) (*Rotator, error) {
	if cfg.Interval <= 0 {
		cfg.Interval = 5 * time.Second
	}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	r := &Rotator{
		path:     path,
		size:     fi.Size(),
		interval: cfg.Interval,
	}
	r.inode = inodeOf(fi)
	return r, nil
}

// Check stats the file and returns ErrRotated if the inode has changed or the
// file has been truncated since the last successful check. On rotation the
// baseline is reset so the next call starts fresh.
func (r *Rotator) Check() error {
	fi, err := os.Stat(r.path)
	if err != nil {
		// File disappeared — treat as rotation.
		return ErrRotated
	}
	curInode := inodeOf(fi)
	curSize := fi.Size()

	if curInode != r.inode || curSize < r.size {
		// Reset baseline for the new file.
		r.inode = curInode
		r.size = curSize
		return ErrRotated
	}
	r.size = curSize
	return nil
}

// Interval returns the configured polling interval.
func (r *Rotator) Interval() time.Duration { return r.interval }
