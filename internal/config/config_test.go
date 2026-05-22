package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logpipe-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `
sources:
  - name: app
    path: /var/log/app.log
    tag: app
sinks:
  - name: console
    type: stdout
workers: 2
`
	path := writeTemp(t, raw)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "app" {
		t.Errorf("expected source name 'app', got %q", cfg.Sources[0].Name)
	}
	if cfg.Workers != 2 {
		t.Errorf("expected workers=2, got %d", cfg.Workers)
	}
}

func TestLoad_DefaultWorkers(t *testing.T) {
	raw := `
sources:
  - name: svc
    path: /tmp/svc.log
sinks:
  - name: out
    type: stdout
`
	cfg, err := config.Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Workers != 4 {
		t.Errorf("expected default workers=4, got %d", cfg.Workers)
	}
}

func TestLoad_MissingSource(t *testing.T) {
	raw := `
sinks:
  - name: out
    type: stdout
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing sources, got nil")
	}
}

func TestLoad_MissingSink(t *testing.T) {
	raw := `
sources:
  - name: app
    path: /tmp/app.log
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing sinks, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
