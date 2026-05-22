package filter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/filter"
)

func TestFilter_Contains(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Contains: "ERROR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("2024/01/01 ERROR something broke") {
		t.Error("expected match for line containing ERROR")
	}
	if f.Match("2024/01/01 INFO all good") {
		t.Error("expected no match for line without ERROR")
	}
}

func TestFilter_Regex(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Regex: `\d{3}-\d{4}`}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("user called 555-1234") {
		t.Error("expected regex match")
	}
	if f.Match("no phone here") {
		t.Error("expected no regex match")
	}
}

func TestFilter_InvalidRegex(t *testing.T) {
	_, err := filter.New([]filter.Rule{{Regex: `[invalid`}})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestFilter_Level(t *testing.T) {
	f, err := filter.New([]filter.Rule{{Level: "warn"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("[WARN] disk usage high") {
		t.Error("expected level match (case-insensitive)")
	}
	if f.Match("[INFO] all clear") {
		t.Error("expected no level match")
	}
}

func TestFilter_MultipleRules(t *testing.T) {
	f, err := filter.New([]filter.Rule{
		{Contains: "ERROR"},
		{Regex: `user=\w+`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("ERROR user=alice failed login") {
		t.Error("expected match when all rules satisfied")
	}
	if f.Match("ERROR no user field") {
		t.Error("expected no match when one rule fails")
	}
}

func TestPassThrough(t *testing.T) {
	f := filter.PassThrough()
	if !f.Match("anything goes") {
		t.Error("pass-through filter should match every line")
	}
	if !f.Match("") {
		t.Error("pass-through filter should match empty line")
	}
}
