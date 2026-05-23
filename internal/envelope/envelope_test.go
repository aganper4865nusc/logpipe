package envelope_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/envelope"
)

func TestNew_FieldsInitialised(t *testing.T) {
	e := envelope.New("src", "hello", 0)
	if e.Fields == nil {
		t.Fatal("expected Fields map to be initialised, got nil")
	}
}

func TestNew_TimestampRecent(t *testing.T) {
	before := time.Now().UTC()
	e := envelope.New("src", "hello", 0)
	after := time.Now().UTC()

	if e.Timestamp.Before(before) || e.Timestamp.After(after) {
		t.Errorf("timestamp %v not between %v and %v", e.Timestamp, before, after)
	}
}

func TestNew_SourceAndLine(t *testing.T) {
	e := envelope.New("mysource", "log line", 42)
	if e.Source != "mysource" {
		t.Errorf("expected source 'mysource', got %q", e.Source)
	}
	if e.Line != "log line" {
		t.Errorf("expected line 'log line', got %q", e.Line)
	}
	if e.Offset != 42 {
		t.Errorf("expected offset 42, got %d", e.Offset)
	}
}

func TestWithField_DoesNotMutateOriginal(t *testing.T) {
	orig := envelope.New("src", "line", 0)
	next := orig.WithField("env", "prod")

	if _, ok := orig.Fields["env"]; ok {
		t.Error("WithField mutated the original envelope")
	}
	if next.Fields["env"] != "prod" {
		t.Errorf("expected field 'env'='prod', got %q", next.Fields["env"])
	}
}

func TestWithField_PreservesExistingFields(t *testing.T) {
	e := envelope.New("src", "line", 0)
	e = e.WithField("a", "1")
	e = e.WithField("b", "2")

	if e.Fields["a"] != "1" {
		t.Errorf("expected a=1, got %q", e.Fields["a"])
	}
	if e.Fields["b"] != "2" {
		t.Errorf("expected b=2, got %q", e.Fields["b"])
	}
}

func TestClone_Independence(t *testing.T) {
	orig := envelope.New("src", "line", 0)
	orig = orig.WithField("k", "v")

	cloned := orig.Clone()
	cloned.Fields["extra"] = "yes"

	if _, ok := orig.Fields["extra"]; ok {
		t.Error("modifying clone affected original")
	}
}
