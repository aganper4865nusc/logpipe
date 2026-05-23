package fanout_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/fanout"
)

type nopSink struct{}

func (nopSink) Write(_ string) error { return nil }

func BenchmarkFanout_2Sinks(b *testing.B) {
	f := fanout.New(nopSink{}, nopSink{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Write("benchmark line")
	}
}

func BenchmarkFanout_8Sinks(b *testing.B) {
	sinks := make([]fanout.Sink, 8)
	for i := range sinks {
		sinks[i] = nopSink{}
	}
	f := fanout.New(sinks...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f.Write("benchmark line")
	}
}
