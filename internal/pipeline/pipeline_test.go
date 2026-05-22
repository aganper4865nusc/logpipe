package pipeline_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/pipeline"
)

// mockSink collects written lines for assertion.
type mockSink struct {
	mu    sync.Mutex
	lines []string
}

func (m *mockSink) Write(line string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lines = append(m.lines, line)
	return nil
}

func TestPipeline_ForwardsLines(t *testing.T) {
	ch := make(chan string, 3)
	ch <- "hello"
	ch <- "world"
	ch <- "foo"
	close(ch)

	s := &mockSink{}
	p := pipeline.New(pipeline.Config{Lines: ch, Sinks: []pipeline.Sink{s}})
	p.Run()

	if len(s.lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(s.lines))
	}
}

func TestPipeline_FilterDropsLines(t *testing.T) {
	f, _ := filter.New([]filter.Rule{{Contains: "ERROR"}})
	ch := make(chan string, 4)
	ch <- "INFO ok"
	ch <- "ERROR bad"
	ch <- "DEBUG trace"
	ch <- "ERROR worse"
	close(ch)

	s := &mockSink{}
	p := pipeline.New(pipeline.Config{Filter: f, Lines: ch, Sinks: []pipeline.Sink{s}})
	p.Run()

	if len(s.lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(s.lines))
	}
	for _, l := range s.lines {
		if !strings.Contains(l, "ERROR") {
			t.Errorf("unexpected line passed filter: %s", l)
		}
	}
}

func TestPipeline_MultipleSinks(t *testing.T) {
	ch := make(chan string, 2)
	ch <- "line1"
	ch <- "line2"
	close(ch)

	s1, s2 := &mockSink{}, &mockSink{}
	p := pipeline.New(pipeline.Config{Lines: ch, Sinks: []pipeline.Sink{s1, s2}})
	p.Run()

	if len(s1.lines) != 2 || len(s2.lines) != 2 {
		t.Errorf("expected both sinks to receive 2 lines, got %d and %d", len(s1.lines), len(s2.lines))
	}
}

func TestPipeline_NilFilterPassesAll(t *testing.T) {
	ch := make(chan string, 1)
	ch <- "anything"
	close(ch)

	s := &mockSink{}
	p := pipeline.New(pipeline.Config{Lines: ch, Sinks: []pipeline.Sink{s}})
	p.Run()

	if len(s.lines) != 1 {
		t.Errorf("expected 1 line with nil filter, got %d", len(s.lines))
	}
}
