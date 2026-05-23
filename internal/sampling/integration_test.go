package sampling_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/sampling"
)

// TestSampler_ConcurrentSafe verifies the sampler is safe for concurrent use.
func TestSampler_ConcurrentSafe(t *testing.T) {
	s, err := sampling.New(sampling.Config{Rate: 0.5, Seed: 7})
	if err != nil {
		t.Fatal(err)
	}

	const goroutines = 20
	const linesEach = 500

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < linesEach; i++ {
				s.Sample("concurrent line")
			}
		}()
	}
	wg.Wait()

	type statter interface {
		Stats() (int64, int64)
	}
	if st, ok := s.(statter); ok {
		total, _ := st.Stats()
		want := int64(goroutines * linesEach)
		if total != want {
			t.Errorf("total lines: want %d, got %d", want, total)
		}
	}
}
