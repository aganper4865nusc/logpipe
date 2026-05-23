// Package offset tracks the read position (byte offset) for each tailed
// source so that logpipe can resume from where it left off after a restart.
// Offsets are stored in-memory and flushed to a checkpoint store on demand.
package offset

import (
	"fmt"
	"sync"
)

// Store holds per-source byte offsets.
type Store struct {
	mu      sync.RWMutex
	offsets map[string]int64
}

// New returns an empty Store.
func New() *Store {
	return &Store{
		offsets: make(map[string]int64),
	}
}

// Set records the current offset for the named source.
func (s *Store) Set(source string, offset int64) error {
	if source == "" {
		return fmt.Errorf("offset: source name must not be empty")
	}
	if offset < 0 {
		return fmt.Errorf("offset: offset must be non-negative, got %d", offset)
	}
	s.mu.Lock()
	s.offsets[source] = offset
	s.mu.Unlock()
	return nil
}

// Get returns the stored offset for source and whether it was found.
func (s *Store) Get(source string) (int64, bool) {
	s.mu.RLock()
	v, ok := s.offsets[source]
	s.mu.RUnlock()
	return v, ok
}

// Delete removes the offset entry for source.
func (s *Store) Delete(source string) {
	s.mu.Lock()
	delete(s.offsets, source)
	s.mu.Unlock()
}

// Sources returns the names of all tracked sources.
func (s *Store) Sources() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.offsets))
	for k := range s.offsets {
		out = append(out, k)
	}
	return out
}

// Snapshot returns a shallow copy of the current offset map.
func (s *Store) Snapshot() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]int64, len(s.offsets))
	for k, v := range s.offsets {
		copy[k] = v
	}
	return copy
}
