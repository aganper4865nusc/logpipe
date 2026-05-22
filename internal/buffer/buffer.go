// Package buffer provides a bounded, thread-safe ring buffer for log lines.
// It is used to decouple fast producers (tailers) from slower sinks.
package buffer

import (
	"errors"
	"sync"
)

// ErrFull is returned when the buffer is at capacity and DropOnFull is false.
var ErrFull = errors.New("buffer: ring buffer is full")

// RingBuffer holds a fixed number of string entries.
type RingBuffer struct {
	mu         sync.Mutex
	items      []string
	head       int
	tail       int
	count      int
	cap        int
	dropOnFull bool
	dropped    int64
}

// New creates a RingBuffer with the given capacity.
// If dropOnFull is true, the oldest entry is silently dropped when full.
func New(capacity int, dropOnFull bool) *RingBuffer {
	if capacity <= 0 {
		capacity = 1
	}
	return &RingBuffer{
		items:      make([]string, capacity),
		cap:        capacity,
		dropOnFull: dropOnFull,
	}
}

// Push adds a line to the buffer. Returns ErrFull if full and dropOnFull is false.
func (r *RingBuffer) Push(line string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == r.cap {
		if !r.dropOnFull {
			return ErrFull
		}
		// Overwrite oldest entry.
		r.head = (r.head + 1) % r.cap
		r.count--
		r.dropped++
	}

	r.items[r.tail] = line
	r.tail = (r.tail + 1) % r.cap
	r.count++
	return nil
}

// Pop removes and returns the oldest entry. Returns ("", false) if empty.
func (r *RingBuffer) Pop() (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return "", false
	}

	line := r.items[r.head]
	r.items[r.head] = ""
	r.head = (r.head + 1) % r.cap
	r.count--
	return line, true
}

// Len returns the current number of buffered entries.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Dropped returns the total number of entries dropped due to overflow.
func (r *RingBuffer) Dropped() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.dropped
}
