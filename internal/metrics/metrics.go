package metrics

import (
	"sync"
	"sync/atomic"
)

// Counters holds runtime statistics for the logpipe process.
type Counters struct {
	LinesRead      atomic.Int64
	LinesFiltered  atomic.Int64
	LinesForwarded atomic.Int64
	SinkErrors     atomic.Int64
}

// Snapshot is a point-in-time copy of Counters.
type Snapshot struct {
	LinesRead      int64
	LinesFiltered  int64
	LinesForwarded int64
	SinkErrors     int64
}

// Registry holds named counter sets.
type Registry struct {
	mu       sync.RWMutex
	sources  map[string]*Counters
	global   Counters
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{
		sources: make(map[string]*Counters),
	}
}

// Source returns (or creates) the Counters for the named source.
func (r *Registry) Source(name string) *Counters {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.sources[name]; ok {
		return c
	}
	c := &Counters{}
	r.sources[name] = c
	return c
}

// Global returns the aggregate Counters.
func (r *Registry) Global() *Counters {
	return &r.global
}

// RecordRead increments read counters for the named source and global.
func (r *Registry) RecordRead(source string) {
	r.Source(source).LinesRead.Add(1)
	r.global.LinesRead.Add(1)
}

// RecordFiltered increments filtered counters.
func (r *Registry) RecordFiltered(source string) {
	r.Source(source).LinesFiltered.Add(1)
	r.global.LinesFiltered.Add(1)
}

// RecordForwarded increments forwarded counters.
func (r *Registry) RecordForwarded(source string) {
	r.Source(source).LinesForwarded.Add(1)
	r.global.LinesForwarded.Add(1)
}

// RecordSinkError increments sink-error counters.
func (r *Registry) RecordSinkError(source string) {
	r.Source(source).SinkErrors.Add(1)
	r.global.SinkErrors.Add(1)
}

// Snapshot returns a point-in-time copy of the global counters.
func (r *Registry) Snapshot() Snapshot {
	return Snapshot{
		LinesRead:      r.global.LinesRead.Load(),
		LinesFiltered:  r.global.LinesFiltered.Load(),
		LinesForwarded: r.global.LinesForwarded.Load(),
		SinkErrors:     r.global.SinkErrors.Load(),
	}
}

// SourceNames returns all registered source names.
func (r *Registry) SourceNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.sources))
	for k := range r.sources {
		names = append(names, k)
	}
	return names
}
