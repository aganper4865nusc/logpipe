package checkpoint_test

import (
	"fmt"
	"testing"

	"github.com/yourusername/logpipe/internal/checkpoint"
)

func BenchmarkSet(b *testing.B) {
	s, err := checkpoint.New(b.TempDir() + "/cp.json")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Set("bench.log", int64(i))
	}
}

func BenchmarkGet(b *testing.B) {
	s, err := checkpoint.New(b.TempDir() + "/cp.json")
	if err != nil {
		b.Fatal(err)
	}
	_ = s.Set("bench.log", 42)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Get("bench.log")
	}
}

func BenchmarkSetMultipleSources(b *testing.B) {
	s, err := checkpoint.New(b.TempDir() + "/cp.json")
	if err != nil {
		b.Fatal(err)
	}
	sources := make([]string, 10)
	for i := range sources {
		sources[i] = fmt.Sprintf("source-%d.log", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Set(sources[i%len(sources)], int64(i))
	}
}
