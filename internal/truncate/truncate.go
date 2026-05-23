// Package truncate provides a line truncation transformer that caps log lines
// at a configurable byte length, appending an optional suffix to indicate
// truncation occurred.
package truncate

import "fmt"

const (
	DefaultMaxBytes = 4096
	DefaultSuffix   = "...[truncated]"
)

// Config holds options for the Truncator.
type Config struct {
	// MaxBytes is the maximum number of bytes allowed per line before truncation.
	// Defaults to DefaultMaxBytes if zero.
	MaxBytes int

	// Suffix is appended to truncated lines. Defaults to DefaultSuffix if empty.
	Suffix string
}

// Truncator is a transform function that shortens lines exceeding MaxBytes.
type Truncator struct {
	maxBytes int
	suffix   string
}

// New returns a new Truncator with the given Config.
// Zero values are replaced with package defaults.
func New(cfg Config) (*Truncator, error) {
	if cfg.MaxBytes < 0 {
		return nil, fmt.Errorf("truncate: MaxBytes must be >= 0, got %d", cfg.MaxBytes)
	}
	if cfg.MaxBytes == 0 {
		cfg.MaxBytes = DefaultMaxBytes
	}
	if cfg.Suffix == "" {
		cfg.Suffix = DefaultSuffix
	}
	return &Truncator{maxBytes: cfg.MaxBytes, suffix: cfg.Suffix}, nil
}

// Transform returns the line unchanged when it fits within MaxBytes.
// Otherwise it returns a prefix of the line (leaving room for the suffix)
// with the suffix appended.
func (t *Truncator) Transform(line string) string {
	if len(line) <= t.maxBytes {
		return line
	}
	avail := t.maxBytes - len(t.suffix)
	if avail <= 0 {
		// Suffix itself is longer than the limit; just return the suffix trimmed.
		if len(t.suffix) > t.maxBytes {
			return t.suffix[:t.maxBytes]
		}
		return t.suffix
	}
	return line[:avail] + t.suffix
}

// AsTransformFunc returns a func(string) string suitable for use with
// transform.Chain and similar helpers.
func (t *Truncator) AsTransformFunc() func(string) string {
	return t.Transform
}
