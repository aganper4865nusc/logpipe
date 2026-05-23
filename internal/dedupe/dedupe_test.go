package dedupe

import (
	"testing"
	"time"
)

func TestDeduplicator_FirstOccurrenceAllowed(t *testing.T) {
	d := New(5 * time.Second)
	if d.IsDuplicate("hello world") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestDeduplicator_SecondOccurrenceBlocked(t *testing.T) {
	d := New(5 * time.Second)
	d.IsDuplicate("repeated line")
	if !d.IsDuplicate("repeated line") {
		t.Fatal("second occurrence within window should be a duplicate")
	}
}

func TestDeduplicator_DifferentLinesAllowed(t *testing.T) {
	d := New(5 * time.Second)
	d.IsDuplicate("line one")
	if d.IsDuplicate("line two") {
		t.Fatal("different line should not be a duplicate")
	}
}

func TestDeduplicator_WindowExpiry(t *testing.T) {
	now := time.Now()
	d := New(2 * time.Second)
	d.now = func() time.Time { return now }

	d.IsDuplicate("expiring line")

	// Advance time beyond the window.
	d.now = func() time.Time { return now.Add(3 * time.Second) }

	if d.IsDuplicate("expiring line") {
		t.Fatal("line should be allowed after window expires")
	}
}

func TestDeduplicator_ZeroWindowDisabled(t *testing.T) {
	d := New(0)
	d.IsDuplicate("same")
	if d.IsDuplicate("same") {
		t.Fatal("zero window should disable deduplication")
	}
}

func TestDeduplicator_Len(t *testing.T) {
	d := New(5 * time.Second)
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	d.IsDuplicate("a") // duplicate, should not increase count
	if got := d.Len(); got != 2 {
		t.Fatalf("expected Len 2, got %d", got)
	}
}

func TestDeduplicator_EvictsExpired(t *testing.T) {
	now := time.Now()
	d := New(1 * time.Second)
	d.now = func() time.Time { return now }

	d.IsDuplicate("old")
	if d.Len() != 1 {
		t.Fatal("expected 1 tracked line")
	}

	d.now = func() time.Time { return now.Add(2 * time.Second) }
	d.IsDuplicate("new") // triggers eviction

	if d.Len() != 1 {
		t.Fatalf("expected 1 after eviction, got %d", d.Len())
	}
}
