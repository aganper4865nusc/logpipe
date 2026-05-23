package parser_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/parser"
)

var sink map[string]any

func BenchmarkJSON(b *testing.B) {
	p := parser.JSON()
	line := `{"level":"info","msg":"benchmark","ts":1234567890}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ := p(line)
		sink = m
	}
}

func BenchmarkRegex(b *testing.B) {
	p, _ := parser.Regex(`(?P<level>\w+) (?P<msg>.+)`)
	line := "INFO benchmark message payload"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ := p(line)
		sink = m
	}
}

func BenchmarkKV(b *testing.B) {
	p := parser.KV(" ", "=")
	line := "level=info msg=benchmark ts=1234567890"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ := p(line)
		sink = m
	}
}

func BenchmarkPassThrough(b *testing.B) {
	p := parser.PassThrough()
	line := "raw log line for benchmarking"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, _ := p(line)
		sink = m
	}
}
