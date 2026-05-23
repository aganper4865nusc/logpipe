// Package checkpoint persists tail offsets so logpipe can resume
// reading from where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Store persists and retrieves per-source byte offsets.
type Store struct {
	mu   sync.Mutex
	path string
	data map[string]int64
}

// New opens (or creates) a checkpoint file at path.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]int64),
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Get returns the last saved offset for source, or 0 if not found.
func (s *Store) Get(source string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[source]
}

// Set updates the offset for source and flushes to disk.
func (s *Store) Set(source string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[source] = offset
	return s.flush()
}

// Sources returns all tracked source names.
func (s *Store) Sources() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Store) flush() error {
	f, err := os.CreateTemp("", "checkpoint-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(s.data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, s.path)
}
