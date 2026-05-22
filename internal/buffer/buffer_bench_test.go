package buffer

import "testing"

// BenchmarkPushPop measures the throughput of sequential push/pop pairs.
func BenchmarkPushPop(b *testing.B) {
	buf := New(256, true)
	line := "2024-01-01T00:00:00Z INFO benchmark log line content here"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Push(line)
		buf.Pop()
	}
}

// BenchmarkPushOnly measures push throughput when the buffer wraps around.
func BenchmarkPushOnly(b *testing.B) {
	buf := New(128, true)
	line := "2024-01-01T00:00:00Z WARN some log message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Push(line)
	}
}

// BenchmarkConcurrentPushPop measures throughput under concurrent load.
func BenchmarkConcurrentPushPop(b *testing.B) {
	buf := New(512, true)
	line := "2024-01-01T00:00:00Z ERROR concurrent benchmark"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf.Push(line)
			buf.Pop()
		}
	})
}
