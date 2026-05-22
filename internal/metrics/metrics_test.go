package metrics

import (
	"testing"
)

func TestRegistry_RecordRead(t *testing.T) {
	r := New()
	r.RecordRead("app.log")
	r.RecordRead("app.log")
	r.RecordRead("sys.log")

	if got := r.Source("app.log").LinesRead.Load(); got != 2 {
		t.Fatalf("expected 2 reads for app.log, got %d", got)
	}
	if got := r.Global().LinesRead.Load(); got != 3 {
		t.Fatalf("expected 3 global reads, got %d", got)
	}
}

func TestRegistry_RecordFiltered(t *testing.T) {
	r := New()
	r.RecordFiltered("app.log")

	if got := r.Source("app.log").LinesFiltered.Load(); got != 1 {
		t.Fatalf("expected 1 filtered, got %d", got)
	}
	if got := r.Global().LinesFiltered.Load(); got != 1 {
		t.Fatalf("expected 1 global filtered, got %d", got)
	}
}

func TestRegistry_RecordForwarded(t *testing.T) {
	r := New()
	r.RecordForwarded("app.log")
	r.RecordForwarded("app.log")

	snap := r.Snapshot()
	if snap.LinesForwarded != 2 {
		t.Fatalf("expected 2 forwarded in snapshot, got %d", snap.LinesForwarded)
	}
}

func TestRegistry_RecordSinkError(t *testing.T) {
	r := New()
	r.RecordSinkError("app.log")

	if got := r.Global().SinkErrors.Load(); got != 1 {
		t.Fatalf("expected 1 sink error, got %d", got)
	}
}

func TestRegistry_SourceNames(t *testing.T) {
	r := New()
	r.RecordRead("a.log")
	r.RecordRead("b.log")

	names := r.SourceNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 source names, got %d", len(names))
	}
}

func TestRegistry_Snapshot_ZeroValues(t *testing.T) {
	r := New()
	snap := r.Snapshot()
	if snap.LinesRead != 0 || snap.LinesFiltered != 0 ||
		snap.LinesForwarded != 0 || snap.SinkErrors != 0 {
		t.Fatal("expected all zero snapshot for new registry")
	}
}

func TestRegistry_MultipleSourcesIsolated(t *testing.T) {
	r := New()
	r.RecordRead("x.log")
	r.RecordFiltered("y.log")

	if r.Source("x.log").LinesFiltered.Load() != 0 {
		t.Fatal("x.log should have 0 filtered lines")
	}
	if r.Source("y.log").LinesRead.Load() != 0 {
		t.Fatal("y.log should have 0 read lines")
	}
}
