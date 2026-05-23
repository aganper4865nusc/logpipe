package fanout_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/yourorg/logpipe/internal/fanout"
)

// captureSink records every line written to it.
type captureSink struct {
	lines []string
	err   error
}

func (c *captureSink) Write(line string) error {
	c.lines = append(c.lines, line)
	return c.err
}

func TestFanout_DeliversToAllSinks(t *testing.T) {
	a, b := &captureSink{}, &captureSink{}
	f := fanout.New(a, b)

	if err := f.Write("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range []*captureSink{a, b} {
		if len(s.lines) != 1 || s.lines[0] != "hello" {
			t.Errorf("sink did not receive line, got %v", s.lines)
		}
	}
}

func TestFanout_AddIncreasesLen(t *testing.T) {
	f := fanout.New()
	if f.Len() != 0 {
		t.Fatalf("expected 0, got %d", f.Len())
	}
	f.Add(&captureSink{})
	if f.Len() != 1 {
		t.Fatalf("expected 1, got %d", f.Len())
	}
}

func TestFanout_CollectsErrors(t *testing.T) {
	errA := errors.New("sink-a failed")
	a := &captureSink{err: errA}
	b := &captureSink{}
	f := fanout.New(a, b)

	err := f.Write("line")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var me *fanout.MultiError
	if !errors.As(err, &me) {
		t.Fatalf("expected *MultiError, got %T", err)
	}
	if len(me.Errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(me.Errs))
	}
	if !errors.Is(me.Errs[0], errA) {
		t.Errorf("unexpected inner error: %v", me.Errs[0])
	}
	// successful sink still received the line
	if len(b.lines) != 1 {
		t.Errorf("successful sink should have received line")
	}
}

func TestFanout_ConcurrentWrites(t *testing.T) {
	var count atomic.Int64
	type countSink struct{}
	f := fanout.New()
	for i := 0; i < 4; i++ {
		f.Add(&countingAdapter{&count})
	}

	for i := 0; i < 100; i++ {
		if err := f.Write("msg"); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}
	if count.Load() != 400 {
		t.Errorf("expected 400 writes, got %d", count.Load())
	}
	_ = countSink{}
}

type countingAdapter struct{ n *atomic.Int64 }

func (c *countingAdapter) Write(_ string) error { c.n.Add(1); return nil }
