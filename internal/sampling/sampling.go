// Package sampling provides probabilistic and rate-based log line sampling
// for reducing volume of high-frequency log sources.
package sampling

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Sampler decides whether a given log line should be forwarded.
type Sampler interface {
	Sample(line string) bool
}

// Config holds configuration for a sampler.
type Config struct {
	// Rate is the fraction of lines to keep, in the range (0, 1].
	// A value of 1.0 keeps all lines; 0.1 keeps ~10%.
	Rate float64

	// Seed is used to initialise the random source. 0 uses a default seed.
	Seed int64
}

// probabilistic sampler keeps each line with probability Rate.
type probabilistic struct {
	rate float64
	mu   sync.Mutex
	rng  *rand.Rand

	total    atomic.Int64
	sampled  atomic.Int64
}

// New returns a Sampler that keeps lines with the given probability.
// Rate must be in (0, 1].
func New(cfg Config) (Sampler, error) {
	if cfg.Rate <= 0 || cfg.Rate > 1 {
		return nil, errors.New("sampling: rate must be in (0, 1]")
	}
	src := rand.NewSource(cfg.Seed)
	return &probabilistic{
		rate: cfg.Rate,
		rng:  rand.New(src), //nolint:gosec
	}, nil
}

// Sample returns true if the line should be forwarded.
func (p *probabilistic) Sample(_ string) bool {
	p.total.Add(1)
	p.mu.Lock()
	v := p.rng.Float64()
	p.mu.Unlock()
	if v < p.rate {
		p.sampled.Add(1)
		return true
	}
	return false
}

// Stats returns total lines seen and lines sampled.
func (p *probabilistic) Stats() (total, sampled int64) {
	return p.total.Load(), p.sampled.Load()
}

// PassThrough is a no-op Sampler that keeps every line.
func PassThrough() Sampler { return passThrough{} }

type passThrough struct{}

func (passThrough) Sample(_ string) bool { return true }
