package redact_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/redact"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := redact.New(map[string]string{"bad": "[invalid("})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_ValidPatterns(t *testing.T) {
	_, err := redact.New(map[string]string{
		"token": `token=[A-Za-z0-9]+`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransform_RedactsMatch(t *testing.T) {
	r, err := redact.New(map[string]string{
		"password": `password=[^\s]+`,
	})
	if err != nil {
		t.Fatal(err)
	}

	input := `user=alice password=s3cr3t action=login`
	out := r.Transform(input)

	if strings.Contains(out, "s3cr3t") {
		t.Errorf("sensitive value still present in output: %q", out)
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Errorf("expected [REDACTED] in output, got: %q", out)
	}
}

func TestTransform_NoMatch(t *testing.T) {
	r, err := redact.New(map[string]string{
		"token": `token=[A-Za-z0-9]+`,
	})
	if err != nil {
		t.Fatal(err)
	}

	input := "user=bob action=logout"
	out := r.Transform(input)
	if out != input {
		t.Errorf("expected unchanged line, got: %q", out)
	}
}

func TestTransform_MultipleRules(t *testing.T) {
	r, err := redact.New(map[string]string{
		"password": `password=[^\s]+`,
		"token":    `token=[A-Za-z0-9]+`,
	})
	if err != nil {
		t.Fatal(err)
	}

	input := `user=alice password=secret token=abc123 action=login`
	out := r.Transform(input)

	for _, sensitive := range []string{"secret", "abc123"} {
		if strings.Contains(out, sensitive) {
			t.Errorf("sensitive value %q still present in output: %q", sensitive, out)
		}
	}
}

func TestPassThrough(t *testing.T) {
	line := "some log line with password=visible"
	if got := redact.PassThrough(line); got != line {
		t.Errorf("PassThrough modified line: got %q", got)
	}
}
