// Package fanout distributes a single log line to multiple downstream channels.
package fanout

import "sync"

// Sink is anything that can receive a log line.
type Sink interface {
	Write(line string) error
}

// Fanout copies each incoming line to every registered sink concurrently.
// Errors from individual sinks are collected and returned together.
type Fanout struct {
	mu    sync.RWMutex
	sinks []Sink
}

// New returns a Fanout pre-loaded with the supplied sinks.
func New(sinks ...Sink) *Fanout {
	return &Fanout{sinks: append([]Sink(nil), sinks...)}
}

// Add registers an additional sink at runtime.
func (f *Fanout) Add(s Sink) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sinks = append(f.sinks, s)
}

// Write delivers line to every sink concurrently and returns a combined error
// if one or more sinks fail.
func (f *Fanout) Write(line string) error {
	f.mu.RLock()
	sinks := append([]Sink(nil), f.sinks...)
	f.mu.RUnlock()

	var (
		mu   sync.Mutex
		errs []error
		wg   sync.WaitGroup
	)

	for _, s := range sinks {
		wg.Add(1)
		go func(s Sink) {
			defer wg.Done()
			if err := s.Write(line); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(s)
	}

	wg.Wait()

	if len(errs) == 0 {
		return nil
	}
	return &MultiError{Errs: errs}
}

// Len returns the current number of registered sinks.
func (f *Fanout) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.sinks)
}
