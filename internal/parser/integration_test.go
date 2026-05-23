package parser_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/parser"
)

// TestParsers_ConcurrentSafe verifies that all parser constructors produce
// ParseFns that are safe to call concurrently from multiple goroutines.
func TestParsers_ConcurrentSafe(t *testing.T) {
	parsers := []struct {
		name string
		fn   func() (parser.ParseFn, error)
		line string
	}{
		{
			name: "json",
			fn:   func() (parser.ParseFn, error) { return parser.JSON(), nil },
			line: `{"level":"info"}`,
		},
		{
			name: "regex",
			fn: func() (parser.ParseFn, error) {
				return parser.Regex(`(?P<level>\w+)`)
			},
			line: "INFO",
		},
		{
			name: "kv",
			fn:   func() (parser.ParseFn, error) { return parser.KV(" ", "="), nil },
			line: "level=info msg=ok",
		},
	}

	const goroutines = 20
	const iterations = 50

	for _, tc := range parsers {
		t.Run(tc.name, func(t *testing.T) {
			parseFn, err := tc.fn()
			if err != nil {
				t.Fatalf("setup: %v", err)
			}
			var wg sync.WaitGroup
			for g := 0; g < goroutines; g++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i < iterations; i++ {
						_, _ = parseFn(tc.line)
					}
				}()
			}
			wg.Wait()
		})
	}
}
