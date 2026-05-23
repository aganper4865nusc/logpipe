package schema_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/schema"
)

// TestValidator_ConcurrentSafe verifies that a single Validator can be used
// from multiple goroutines without data races.
func TestValidator_ConcurrentSafe(t *testing.T) {
	v := schema.New([]string{"level", "msg"})

	lines := []string{
		`{"level":"info","msg":"ok"}`,
		`{"level":"warn"}`,
		`not json`,
		`{"level":"error","msg":"boom"}`,
	}

	var wg sync.WaitGroup
	const workers = 20
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for _, l := range lines {
				_ = v.Validate(l) // result intentionally ignored; race detector is the goal
			}
		}(i)
	}

	wg.Wait()
}
