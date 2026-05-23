package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/logpipe/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_CreatesEmptyStore(t *testing.T) {
	s, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("source1"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestSet_PersistsOffset(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	if err := s.Set("app.log", 1024); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Reload from disk.
	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if got := s2.Get("app.log"); got != 1024 {
		t.Errorf("expected 1024, got %d", got)
	}
}

func TestSet_MultipleSourcesIndependent(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("a.log", 100)
	_ = s.Set("b.log", 200)

	if s.Get("a.log") != 100 {
		t.Errorf("a.log: expected 100")
	}
	if s.Get("b.log") != 200 {
		t.Errorf("b.log: expected 200")
	}
}

func TestSources_ReturnsAllKeys(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Set("x.log", 1)
	_ = s.Set("y.log", 2)

	srcs := s.Sources()
	if len(srcs) != 2 {
		t.Errorf("expected 2 sources, got %d", len(srcs))
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	p := tempPath(t)
	if err := os.WriteFile(p, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := checkpoint.New(p); err == nil {
		t.Error("expected error for corrupt checkpoint file")
	}
}

func TestGet_UnknownSource_ReturnsZero(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	if got := s.Get("missing"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}
