package masker

import (
	"encoding/json"
	"testing"
)

func TestNew_DefaultMask(t *testing.T) {
	m := New([]string{"password"}, "")
	if m.mask != defaultMask {
		t.Fatalf("expected default mask %q, got %q", defaultMask, m.mask)
	}
}

func TestTransform_NonJSON_Unchanged(t *testing.T) {
	m := New([]string{"password"}, "")
	line := "plain text log line"
	if got := m.Transform(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestTransform_MasksField(t *testing.T) {
	m := New([]string{"password"}, "REDACTED")
	line := `{"user":"alice","password":"s3cr3t"}`
	out := m.Transform(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["password"] != "REDACTED" {
		t.Errorf("expected password to be REDACTED, got %v", obj["password"])
	}
	if obj["user"] != "alice" {
		t.Errorf("expected user to be unchanged, got %v", obj["user"])
	}
}

func TestTransform_MultipleFields(t *testing.T) {
	m := New([]string{"token", "secret"}, "")
	line := `{"token":"abc","secret":"xyz","level":"info"}`
	out := m.Transform(line)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["token"] != defaultMask {
		t.Errorf("token not masked: %v", obj["token"])
	}
	if obj["secret"] != defaultMask {
		t.Errorf("secret not masked: %v", obj["secret"])
	}
	if obj["level"] != "info" {
		t.Errorf("level should be unchanged: %v", obj["level"])
	}
}

func TestTransform_NoMatchingField(t *testing.T) {
	m := New([]string{"password"}, "")
	line := `{"user":"bob","level":"warn"}`
	out := m.Transform(line)

	var orig, got map[string]interface{}
	_ = json.Unmarshal([]byte(line), &orig)
	_ = json.Unmarshal([]byte(out), &got)
	if got["user"] != orig["user"] || got["level"] != orig["level"] {
		t.Errorf("unexpected mutation: %v", got)
	}
}

func TestTransform_EmptyFields_NoOp(t *testing.T) {
	m := New(nil, "")
	line := `{"password":"secret"}`
	if got := m.Transform(line); got != line {
		t.Fatalf("expected no-op, got %q", got)
	}
}

func TestPassThrough(t *testing.T) {
	line := `{"password":"secret"}`
	if got := PassThrough(line); got != line {
		t.Fatalf("PassThrough mutated line: %q", got)
	}
}
