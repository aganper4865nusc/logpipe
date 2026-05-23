package aggregate

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNew_NilFunc(t *testing.T) {
	_, err := New(10, 0, nil)
	if err == nil {
		t.Fatal("expected error for nil flush func")
	}
}

func TestNew_NoFlushCondition(t *testing.T) {
	_, err := New(0, 0, func([]string) {})
	if err == nil {
		t.Fatal("expected error when both maxSize and interval are zero")
	}
}

func TestAdd_FlushesOnSize(t *testing.T) {
	var mu sync.Mutex
	var got []string
	agg, err := New(3, 0, func(lines []string) {
		mu.Lock()
		got = append(got, lines...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}

	// Run with a context that won't cancel during the test.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go agg.Run(ctx)

	agg.Add("a")
	agg.Add("b")
	agg.Add("c") // triggers flush

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 3 {
		t.Fatalf("expected 3 lines flushed, got %d", len(got))
	}
}

func TestRun_FlushesOnInterval(t *testing.T) {
	var mu sync.Mutex
	var flushed [][]string
	agg, err := New(0, 30*time.Millisecond, func(lines []string) {
		mu.Lock()
		flushed = append(flushed, lines)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go agg.Run(ctx)

	agg.Add("x")
	agg.Add("y")
	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(flushed) == 0 {
		t.Fatal("expected at least one interval flush")
	}
	if len(flushed[0]) != 2 {
		t.Fatalf("expected 2 lines in first flush, got %d", len(flushed[0]))
	}
}

func TestRun_FinalFlushOnCancel(t *testing.T) {
	var mu sync.Mutex
	var got []string
	agg, err := New(0, time.Hour, func(lines []string) {
		mu.Lock()
		got = append(got, lines...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	agg.Add("final")

	done := make(chan struct{})
	go func() {
		agg.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancel")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || got[0] != "final" {
		t.Fatalf("expected final flush with [final], got %v", got)
	}
}

func TestFlush_EmptyBufferIsNoop(t *testing.T) {
	called := false
	agg, _ := New(10, 0, func(lines []string) { called = true })
	agg.Flush()
	if called {
		t.Fatal("flush func should not be called on empty buffer")
	}
}
