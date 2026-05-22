package sink

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStdoutSink_Write(t *testing.T) {
	var buf bytes.Buffer
	s := &StdoutSink{w: &buf}

	if err := s.Write("hello stdout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello stdout" {
		t.Errorf("expected %q, got %q", "hello stdout", got)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
}

func TestFileSink_Write(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	s, err := NewFileSink(path)
	if err != nil {
		t.Fatalf("NewFileSink: %v", err)
	}
	defer s.Close()

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := s.Write(l); err != nil {
			t.Fatalf("Write(%q): %v", l, err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	for _, l := range lines {
		if !strings.Contains(string(data), l) {
			t.Errorf("expected file to contain %q", l)
		}
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("kafka", "")
	if err == nil {
		t.Fatal("expected error for unknown sink type")
	}
}

func TestNew_FileMissingTarget(t *testing.T) {
	_, err := New("file", "")
	if err == nil {
		t.Fatal("expected error for missing file target")
	}
}

func TestNew_Stdout(t *testing.T) {
	s, err := New("stdout", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.(*StdoutSink); !ok {
		t.Errorf("expected *StdoutSink, got %T", s)
	}
}
