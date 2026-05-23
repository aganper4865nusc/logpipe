package schema_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/schema"
)

var sink error

func BenchmarkValidate_Valid(b *testing.B) {
	v := schema.New([]string{"level", "msg", "ts"})
	line := `{"level":"info","msg":"benchmark","ts":"2024-01-01T00:00:00Z"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = v.Validate(line)
	}
}

func BenchmarkValidate_Missing(b *testing.B) {
	v := schema.New([]string{"level", "msg", "ts"})
	line := `{"level":"info"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = v.Validate(line)
	}
}

func BenchmarkPassThrough(b *testing.B) {
	v := schema.PassThrough()
	line := `{"level":"info","msg":"benchmark"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = v.Validate(line)
	}
}
