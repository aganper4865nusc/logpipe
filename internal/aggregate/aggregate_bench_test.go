package aggregate

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func BenchmarkAdd_SizeFlush(b *testing.B) {
	agg, _ := New(64, 0, func(lines []string) {})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go agg.Run(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agg.Add(strconv.Itoa(i))
	}
}

func BenchmarkAdd_IntervalFlush(b *testing.B) {
	agg, _ := New(0, 100*time.Millisecond, func(lines []string) {})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go agg.Run(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agg.Add(strconv.Itoa(i))
	}
}
