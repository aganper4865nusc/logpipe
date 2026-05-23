package sampling_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/sampling"
)

func BenchmarkSample_Rate1(b *testing.B) {
	s, _ := sampling.New(sampling.Config{Rate: 1.0, Seed: 1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Sample("benchmark log line")
	}
}

func BenchmarkSample_Rate01(b *testing.B) {
	s, _ := sampling.New(sampling.Config{Rate: 0.1, Seed: 2})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Sample("benchmark log line")
	}
}

func BenchmarkPassThrough(b *testing.B) {
	s := sampling.PassThrough()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Sample("benchmark log line")
	}
}
