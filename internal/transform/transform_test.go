package transform_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/transform"
)

func TestPassThrough(t *testing.T) {
	t.Parallel()
	tr := transform.PassThrough()
	out, err := tr("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", out)
	}
}

func TestUppercase(t *testing.T) {
	t.Parallel()
	tr := transform.Uppercase()
	out, err := tr("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "HELLO" {
		t.Fatalf("expected HELLO, got %q", out)
	}
}

func TestAddField_NonJSON(t *testing.T) {
	t.Parallel()
	tr := transform.AddField("source", "app")
	out, err := tr("plain text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["source"] != "app" {
		t.Fatalf("expected source=app, got %v", obj["source"])
	}
	if obj["message"] != "plain text" {
		t.Fatalf("expected message=plain text, got %v", obj["message"])
	}
}

func TestAddField_ExistingJSON(t *testing.T) {
	t.Parallel()
	tr := transform.AddField("env", "prod")
	out, err := tr(`{"level":"info"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", obj["env"])
	}
	if obj["level"] != "info" {
		t.Fatalf("level field missing")
	}
}

func TestAddTimestamp_InjectsKey(t *testing.T) {
	t.Parallel()
	tr := transform.AddTimestamp("ts")
	out, err := tr(`{"msg":"hi"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["ts"]; !ok {
		t.Fatal("expected ts key in output")
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	t.Parallel()
	tr := transform.Chain(
		transform.AddField("app", "logpipe"),
		transform.AddField("version", "1"),
	)
	out, err := tr(`{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "logpipe") || !strings.Contains(out, "\"version\"") {
		t.Fatalf("chain did not apply all transformers: %s", out)
	}
}

func TestChain_StopsOnError(t *testing.T) {
	t.Parallel()
	failTr := func(_ string) (string, error) {
		return "", &testError{"boom"}
	}
	tr := transform.Chain(failTr, transform.Uppercase())
	_, err := tr("hello")
	if err == nil {
		t.Fatal("expected error from chain")
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
