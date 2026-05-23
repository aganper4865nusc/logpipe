// Package dedupe provides a line deduplication filter that suppresses
// repeated log lines within a configurable time window.
package dedupe

import (
	"sync"
	"time"
)

// Deduplicator tracks recently seen lines and suppresses duplicates
// within a sliding window.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// New returns a Deduplicator that suppresses duplicate lines seen within
// the given window duration. Pass a zero duration to disable deduplication.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// IsDuplicate reports whether line was already seen within the window.
// If it is not a duplicate the line is recorded and false is returned.
func (d *Deduplicator) IsDuplicate(line string) bool {
	if d.window == 0 {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	d.evict(now)

	if _, exists := d.seen[line]; exists {
		return true
	}

	d.seen[line] = now
	return false
}

// evict removes entries whose window has expired. Must be called with mu held.
func (d *Deduplicator) evict(now time.Time) {
	cutoff := now.Add(-d.window)
	for k, t := range d.seen {
		if t.Before(cutoff) {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of currently tracked unique lines.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
