package tail

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()

	out := make(chan Line, 8)
	tlr := New(f.Name(), out)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- tlr.Run(ctx)
	}()

	// Give the tailer time to seek to end before writing.
	time.Sleep(50 * time.Millisecond)

	_, err = f.WriteString("hello world\n")
	if err != nil {
		t.Fatalf("write to file: %v", err)
	}

	select {
	case line := <-out:
		if line.Source != f.Name() {
			t.Errorf("expected source %q, got %q", f.Name(), line.Source)
		}
		if line.Text != "hello world\n" {
			t.Errorf("expected text %q, got %q", "hello world\n", line.Text)
		}
		if line.Time.IsZero() {
			t.Error("expected non-zero timestamp")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for line")
	}

	cancel()
	if runErr := <-errCh; runErr != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", runErr)
	}
}

func TestTailer_MissingFile(t *testing.T) {
	out := make(chan Line, 1)
	tlr := New("/nonexistent/path/file.log", out)

	ctx := context.Background()
	err := tlr.Run(ctx)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
