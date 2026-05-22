package ratelimit_test

import (
	"context"
	"testing"

	"github.com/yourorg/logpipe/internal/ratelimit"
)

func BenchmarkTryConsume(b *testing.B) {
	l := ratelimit.New(ratelimit.Config{Rate: 1e9, Burst: 1e9})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.TryConsume()
	}
}

func BenchmarkWait(b *testing.B) {
	l := ratelimit.New(ratelimit.Config{Rate: 1e9, Burst: 1e9})
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Wait(ctx)
	}
}
