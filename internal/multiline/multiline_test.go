package multiline

import (
	"strings"
	"testing"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(Config{StartPattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidPattern(t *testing.T) {
	_, err := New(Config{StartPattern: `^\d{4}-\d{2}-\d{2}`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFeed_SingleLine(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\[`})
	event, ready := a.Feed("[INFO] hello")
	if ready {
		t.Fatalf("expected not ready on first line, got event %q", event)
	}
}

func TestFeed_ContinuationFolded(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\[`})
	a.Feed("[INFO] start")
	a.Feed("  continuation line")

	event := a.Flush()
	if !strings.Contains(event, "continuation line") {
		t.Errorf("expected continuation in flushed event, got %q", event)
	}
	if !strings.Contains(event, "[INFO] start") {
		t.Errorf("expected start line in flushed event, got %q", event)
	}
}

func TestFeed_NewStartFlushesBuffer(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\[`})
	a.Feed("[INFO] first")
	a.Feed("  detail")

	event, ready := a.Feed("[WARN] second")
	if !ready {
		t.Fatal("expected ready when new start line encountered")
	}
	if !strings.Contains(event, "[INFO] first") {
		t.Errorf("flushed event should contain first line, got %q", event)
	}
	if strings.Contains(event, "[WARN] second") {
		t.Errorf("flushed event should not contain second start line, got %q", event)
	}
}

func TestFeed_MaxLinesFlushes(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\[`, MaxLines: 3})
	a.Feed("[INFO] line1")
	a.Feed("  line2")
	event, ready := a.Feed("  line3")
	if !ready {
		t.Fatal("expected flush at MaxLines")
	}
	parts := strings.Split(event, "\n")
	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d: %q", len(parts), event)
	}
}

func TestFlush_EmptyBuffer(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\[`})
	if got := a.Flush(); got != "" {
		t.Errorf("expected empty string from empty flush, got %q", got)
	}
}

func TestFeed_SequentialEvents(t *testing.T) {
	a, _ := New(Config{StartPattern: `^\d`})
	events := []string{}

	lines := []string{"2024-01-01 first", "  trace A", "2024-01-02 second", "  trace B"}
	for _, l := range lines {
		if ev, ok := a.Feed(l); ok {
			events = append(events, ev)
		}
	}
	events = append(events, a.Flush())

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d: %v", len(events), events)
	}
	if !strings.Contains(events[0], "first") {
		t.Errorf("first event should contain 'first': %q", events[0])
	}
	if !strings.Contains(events[1], "second") {
		t.Errorf("second event should contain 'second': %q", events[1])
	}
}
