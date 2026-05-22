package buffer

import (
	"sync"
	"testing"
)

func TestRingBuffer_PushPop(t *testing.T) {
	b := New(3, false)

	if err := b.Push("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Push("b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line, ok := b.Pop()
	if !ok || line != "a" {
		t.Fatalf("expected 'a', got %q ok=%v", line, ok)
	}
	if b.Len() != 1 {
		t.Fatalf("expected len 1, got %d", b.Len())
	}
}

func TestRingBuffer_FullReturnsError(t *testing.T) {
	b := New(2, false)
	b.Push("x")
	b.Push("y")

	if err := b.Push("z"); err != ErrFull {
		t.Fatalf("expected ErrFull, got %v", err)
	}
}

func TestRingBuffer_DropOnFull(t *testing.T) {
	b := New(2, true)
	b.Push("first")
	b.Push("second")
	b.Push("third") // should drop "first"

	if b.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", b.Dropped())
	}

	line, _ := b.Pop()
	if line != "second" {
		t.Fatalf("expected 'second' after drop, got %q", line)
	}
}

func TestRingBuffer_PopEmpty(t *testing.T) {
	b := New(4, false)
	_, ok := b.Pop()
	if ok {
		t.Fatal("expected ok=false on empty buffer")
	}
}

func TestRingBuffer_WrapAround(t *testing.T) {
	b := New(3, false)
	b.Push("1")
	b.Push("2")
	b.Push("3")
	b.Pop()
	b.Push("4")

	expected := []string{"2", "3", "4"}
	for _, want := range expected {
		got, ok := b.Pop()
		if !ok || got != want {
			t.Fatalf("expected %q, got %q ok=%v", want, got, ok)
		}
	}
}

func TestRingBuffer_ConcurrentAccess(t *testing.T) {
	b := New(64, true)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Push("line")
			b.Pop()
		}()
	}
	wg.Wait()
}
