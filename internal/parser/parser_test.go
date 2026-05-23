package parser_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/parser"
)

func TestJSON_ValidLine(t *testing.T) {
	p := parser.JSON()
	m, err := p(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestJSON_InvalidLine(t *testing.T) {
	p := parser.JSON()
	_, err := p("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRegex_NamedGroups(t *testing.T) {
	p, err := parser.Regex(`(?P<level>\w+) (?P<msg>.+)`)
	if err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	m, err := p("INFO server started")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "INFO" || m["msg"] != "server started" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestRegex_NoMatch(t *testing.T) {
	p, _ := parser.Regex(`^\d+$`)
	_, err := p("not-a-number")
	if err == nil {
		t.Fatal("expected error when line does not match")
	}
}

func TestRegex_InvalidPattern(t *testing.T) {
	_, err := parser.Regex(`[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestKV_BasicPairs(t *testing.T) {
	p := parser.KV(" ", "=")
	m, err := p("level=info msg=hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestKV_NoPairsReturnsError(t *testing.T) {
	p := parser.KV(" ", "=")
	_, err := p("no-pairs-here")
	if err == nil {
		t.Fatal("expected error when no pairs found")
	}
}

func TestPassThrough_WrapsMessage(t *testing.T) {
	p := parser.PassThrough()
	m, err := p("raw log line")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["message"] != "raw log line" {
		t.Fatalf("expected message key, got: %v", m)
	}
}
