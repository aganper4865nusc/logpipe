package rotator

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestNew_DefaultInterval(t *testing.T) {
	p := tempFile(t, "hello\n")
	r, err := New(p, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if r.Interval() != 5*time.Second {
		t.Fatalf("expected 5s default, got %v", r.Interval())
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := New("/no/such/file.log", Config{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheck_NoRotation(t *testing.T) {
	p := tempFile(t, "line1\n")
	r, err := New(p, Config{Interval: time.Second})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Append data — size grows, inode unchanged.
	f, _ := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0o644)
	_, _ = f.WriteString("line2\n")
	f.Close()

	if err := r.Check(); err != nil {
		t.Fatalf("unexpected rotation signal: %v", err)
	}
}

func TestCheck_TruncationDetected(t *testing.T) {
	p := tempFile(t, "some content\n")
	r, err := New(p, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Truncate the file.
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	if err := r.Check(); err != ErrRotated {
		t.Fatalf("expected ErrRotated, got %v", err)
	}
}

func TestCheck_FileDisappears(t *testing.T) {
	p := tempFile(t, "data\n")
	r, err := New(p, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	os.Remove(p)
	if err := r.Check(); err != ErrRotated {
		t.Fatalf("expected ErrRotated, got %v", err)
	}
}

func TestCheck_ResetAfterRotation(t *testing.T) {
	p := tempFile(t, "initial\n")
	r, err := New(p, Config{})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Truncate → rotation detected.
	_ = os.WriteFile(p, []byte(""), 0o644)
	_ = r.Check()
	// Append new data; no further rotation should be reported.
	f, _ := os.OpenFile(p, os.O_APPEND|os.O_WRONLY, 0o644)
	_, _ = f.WriteString("new\n")
	f.Close()
	if err := r.Check(); err != nil {
		t.Fatalf("unexpected second rotation: %v", err)
	}
}
