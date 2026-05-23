// Package splitter provides a fan-out mechanism that duplicates each log line
// to multiple independent output channels, enabling parallel processing paths.
package splitter

import "sync"

// Splitter fans a single stream of lines out to N output channels.
type Splitter struct {
	mu      sync.RWMutex
	outputs []chan string
}

// New creates a Splitter with the given number of output channels.
// Each channel is buffered to size bufSize.
func New(n, bufSize int) (*Splitter, []<-chan string) {
	if n <= 0 {
		n = 1
	}
	if bufSize <= 0 {
		bufSize = 64
	}

	s := &Splitter{
		outputs: make([]chan string, n),
	}
	readers := make([]<-chan string, n)
	for i := 0; i < n; i++ {
		ch := make(chan string, bufSize)
		s.outputs[i] = ch
		readers[i] = ch
	}
	return s, readers
}

// Write sends line to every registered output channel.
// If a channel is full the line is dropped for that output only.
func (s *Splitter) Write(line string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.outputs {
		select {
		case ch <- line:
		default:
			// drop rather than block
		}
	}
}

// Close closes all output channels, signalling consumers that no more lines
// will arrive.
func (s *Splitter) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.outputs {
		close(ch)
	}
}

// Len returns the number of output channels.
func (s *Splitter) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.outputs)
}
