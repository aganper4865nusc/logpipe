// Package routing provides log-line routing based on source tags and rules.
// A Router maps named sources to one or more named sinks, allowing selective
// forwarding of log lines without duplicating pipeline configuration.
package routing

import (
	"fmt"
	"sync"
)

// Route defines a mapping from a source name to a set of sink names.
type Route struct {
	Source string
	Sinks  []string
}

// Router holds a set of routes and resolves sink names for a given source.
type Router struct {
	mu     sync.RWMutex
	routes map[string][]string // source -> []sink
}

// New creates a Router from the provided routes.
// Returns an error if any route has an empty source or no sinks.
func New(routes []Route) (*Router, error) {
	r := &Router{
		routes: make(map[string][]string, len(routes)),
	}
	for _, rt := range routes {
		if rt.Source == "" {
			return nil, fmt.Errorf("routing: route has empty source")
		}
		if len(rt.Sinks) == 0 {
			return nil, fmt.Errorf("routing: source %q has no sinks", rt.Source)
		}
		r.routes[rt.Source] = append(r.routes[rt.Source], rt.Sinks...)
	}
	return r, nil
}

// Resolve returns the sink names associated with the given source.
// If no route is registered for the source, it returns nil, false.
func (r *Router) Resolve(source string) ([]string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sinks, ok := r.routes[source]
	if !ok {
		return nil, false
	}
	out := make([]string, len(sinks))
	copy(out, sinks)
	return out, true
}

// Add registers or appends a route at runtime.
func (r *Router) Add(source string, sinks []string) error {
	if source == "" {
		return fmt.Errorf("routing: source must not be empty")
	}
	if len(sinks) == 0 {
		return fmt.Errorf("routing: must provide at least one sink for source %q", source)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes[source] = append(r.routes[source], sinks...)
	return nil
}

// Sources returns all registered source names.
func (r *Router) Sources() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.routes))
	for k := range r.routes {
		names = append(names, k)
	}
	return names
}
