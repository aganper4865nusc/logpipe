package label

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNew_EmptyLabels(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty labels")
	}
}

func TestNew_NilLabels(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil labels")
	}
}

func TestNew_BlankKey(t *testing.T) {
	_, err := New(map[string]string{"  ": "value"})
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestTransform_InjectsIntoJSON(t *testing.T) {
	l, err := New(map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := l.Transform(`{"msg":"hello"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", obj["env"])
	}
	if obj["msg"] != "hello" {
		t.Errorf("original field lost, got %v", obj["msg"])
	}
}

func TestTransform_DoesNotOverwriteExistingJSONKey(t *testing.T) {
	l, _ := New(map[string]string{"env": "prod"})
	out := l.Transform(`{"env":"staging"}`)
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["env"] != "staging" {
		t.Errorf("existing key should not be overwritten, got %v", obj["env"])
	}
}

func TestTransform_PlainTextAppendsPairs(t *testing.T) {
	l, _ := New(map[string]string{"host": "box1"})
	out := l.Transform("some plain log line")
	if !strings.HasPrefix(out, "some plain log line") {
		t.Errorf("original text missing from output: %q", out)
	}
	if !strings.Contains(out, "host=box1") {
		t.Errorf("expected host=box1 in output: %q", out)
	}
}

func TestPassThrough_Identity(t *testing.T) {
	line := "unchanged line"
	if got := PassThrough(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}
