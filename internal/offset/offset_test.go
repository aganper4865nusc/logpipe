package offset

import (
	"sort"
	"testing"
)

func TestNew_EmptyStore(t *testing.T) {
	s := New()
	if len(s.Sources()) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(s.Sources()))
	}
}

func TestSet_RecordsOffset(t *testing.T) {
	s := New()
	if err := s.Set("app.log", 1024); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("app.log")
	if !ok {
		t.Fatal("expected offset to be present")
	}
	if v != 1024 {
		t.Fatalf("expected 1024, got %d", v)
	}
}

func TestSet_EmptySourceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", 0); err == nil {
		t.Fatal("expected error for empty source name")
	}
}

func TestSet_NegativeOffsetReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("app.log", -1); err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestGet_UnknownSourceReturnsFalse(t *testing.T) {
	s := New()
	_, ok := s.Get("missing.log")
	if ok {
		t.Fatal("expected ok=false for unknown source")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("a.log", 42)
	s.Delete("a.log")
	_, ok := s.Get("a.log")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestSources_ReturnsAllKeys(t *testing.T) {
	s := New()
	_ = s.Set("x.log", 1)
	_ = s.Set("y.log", 2)
	_ = s.Set("z.log", 3)
	srcs := s.Sources()
	sort.Strings(srcs)
	if len(srcs) != 3 || srcs[0] != "x.log" || srcs[1] != "y.log" || srcs[2] != "z.log" {
		t.Fatalf("unexpected sources: %v", srcs)
	}
}

func TestSnapshot_IsIndependentCopy(t *testing.T) {
	s := New()
	_ = s.Set("a.log", 100)
	snap := s.Snapshot()
	_ = s.Set("a.log", 999)
	if snap["a.log"] != 100 {
		t.Fatalf("snapshot should not reflect later writes, got %d", snap["a.log"])
	}
}
