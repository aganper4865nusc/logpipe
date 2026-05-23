package schema_test

import (
	"errors"
	"testing"

	"github.com/yourorg/logpipe/internal/schema"
)

func TestNew_NoRequired_AcceptsAnything(t *testing.T) {
	v := schema.New(nil)
	if err := v.Validate("not json at all"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestPassThrough_AcceptsAnything(t *testing.T) {
	v := schema.PassThrough()
	if err := v.Validate("{}"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestValidate_ValidJSON_AllFields(t *testing.T) {
	v := schema.New([]string{"level", "msg"})
	line := `{"level":"info","msg":"hello"}`
	if err := v.Validate(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_MissingField(t *testing.T) {
	v := schema.New([]string{"level", "msg", "ts"})
	line := `{"level":"info","msg":"hello"}`
	err := v.Validate(line)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *schema.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(ve.Missing) != 1 || ve.Missing[0] != "ts" {
		t.Fatalf("unexpected missing fields: %v", ve.Missing)
	}
}

func TestValidate_NonJSON_WithRequired(t *testing.T) {
	v := schema.New([]string{"level"})
	err := v.Validate("plain text log line")
	if err == nil {
		t.Fatal("expected error for non-JSON line")
	}
}

func TestValidate_MultipleMissingFields(t *testing.T) {
	v := schema.New([]string{"a", "b", "c"})
	err := v.Validate(`{"a":1}`)
	var ve *schema.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if len(ve.Missing) != 2 {
		t.Fatalf("expected 2 missing, got %d: %v", len(ve.Missing), ve.Missing)
	}
}

func TestValidationError_Message(t *testing.T) {
	ve := &schema.ValidationError{Missing: []string{"x", "y"}}
	msg := ve.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}
