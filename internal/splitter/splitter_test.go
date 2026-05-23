package splitter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/splitter"
)

func drain(ch <-chan string, timeout time.Duration) []string {
	var out []string
	deadline := time.After(timeout)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-deadline:
			return out
		}
	}
}

func TestSplitter_Len(t *testing.T) {
	s, _ := splitter.New(3, 8)
	if got := s.Len(); got != 3 {
		t.Fatalf("expected 3 outputs, got %d", got)
	}
}

func TestSplitter_DefaultsApplied(t *testing.T) {
	s, readers := splitter.New(0, 0)
	if s.Len() != 1 {
		t.Fatalf("expected 1 output for n=0, got %d", s.Len())
	}
	if len(readers) != 1 {
		t.Fatalf("expected 1 reader, got %d", len(readers))
	}
}

func TestSplitter_WriteFanOut(t *testing.T) {
	s, readers := splitter.New(3, 16)

	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		s.Write(l)
	}
	s.Close()

	for i, r := range readers {
		got := drain(r, time.Second)
		if len(got) != len(lines) {
			t.Fatalf("reader %d: expected %d lines, got %d", i, len(lines), len(got))
		}
		for j, want := range lines {
			if got[j] != want {
				t.Errorf("reader %d line %d: want %q, got %q", i, j, want, got[j])
			}
		}
	}
}

func TestSplitter_DropWhenFull(t *testing.T) {
	// bufSize=1 so the second write to a slow consumer is dropped
	s, readers := splitter.New(1, 1)

	s.Write("first")
	s.Write("dropped") // channel full, should not block
	s.Close()

	got := drain(readers[0], time.Second)
	if len(got) != 1 {
		t.Fatalf("expected 1 line (drop on full), got %d: %v", len(got), got)
	}
	if got[0] != "first" {
		t.Errorf("expected %q, got %q", "first", got[0])
	}
}

func TestSplitter_CloseSignalsConsumers(t *testing.T) {
	s, readers := splitter.New(2, 8)
	s.Write("x")
	s.Close()

	for i, r := range readers {
		select {
		case _, ok := <-r:
			_ = ok
		case <-time.After(time.Second):
			t.Errorf("reader %d: channel not closed after Close()", i)
		}
	}
}
