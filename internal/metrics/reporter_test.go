package metrics

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestReporter_WritesSnapshot(t *testing.T) {
	reg := New()
	reg.RecordRead("a.log")
	reg.RecordRead("a.log")
	reg.RecordForwarded("a.log")
	reg.RecordFiltered("a.log")
	reg.RecordSinkError("a.log")

	var buf bytes.Buffer
	reporter := NewReporter(reg, &buf, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	reporter.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "read=2") {
		t.Errorf("expected read=2 in output, got: %s", out)
	}
	if !strings.Contains(out, "filtered=1") {
		t.Errorf("expected filtered=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "forwarded=1") {
		t.Errorf("expected forwarded=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "sink_errors=1") {
		t.Errorf("expected sink_errors=1 in output, got: %s", out)
	}
}

func TestReporter_TicksMultipleTimes(t *testing.T) {
	reg := New()
	var buf bytes.Buffer
	reporter := NewReporter(reg, &buf, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Millisecond)
	defer cancel()

	reporter.Run(ctx)

	lines := strings.Count(buf.String(), "[metrics]")
	if lines < 2 {
		t.Errorf("expected at least 2 metric lines, got %d", lines)
	}
}
