package batch

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBatcher_FlushOnSize(t *testing.T) {
	var mu sync.Mutex
	var got [][]string

	b := New(3, time.Second, func(lines []string) {
		mu.Lock()
		got = append(got, lines)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go b.Run(ctx)

	b.Send("a")
	b.Send("b")
	b.Send("c")

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected at least one flush")
	}
	if len(got[0]) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got[0]))
	}
}

func TestBatcher_FlushOnTimeout(t *testing.T) {
	var mu sync.Mutex
	var got [][]string

	b := New(100, 50*time.Millisecond, func(lines []string) {
		mu.Lock()
		got = append(got, lines)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go b.Run(ctx)

	b.Send("x")
	b.Send("y")

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected timeout flush")
	}
	if len(got[0]) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got[0]))
	}
}

func TestBatcher_FlushOnContextCancel(t *testing.T) {
	var mu sync.Mutex
	var got [][]string

	b := New(100, time.Second, func(lines []string) {
		mu.Lock()
		got = append(got, lines)
		mu.Unlock()
	})

	ctx, cancel := context.WithCancel(context.Background())
	go b.Run(ctx)

	b.Send("line1")
	time.Sleep(20 * time.Millisecond)
	cancel()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected flush on cancel")
	}
}

func TestBatcher_DefaultsApplied(t *testing.T) {
	b := New(0, 0, func(_ []string) {})
	if b.maxSize != 100 {
		t.Errorf("expected default maxSize 100, got %d", b.maxSize)
	}
	if b.timeout != time.Second {
		t.Errorf("expected default timeout 1s, got %v", b.timeout)
	}
}
