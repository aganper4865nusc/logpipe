package format_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/format"
)

func TestPassThrough(t *testing.T) {
	line := "hello world"
	if got := format.PassThrough(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestJSON_PlainLine(t *testing.T) {
	fmt := format.JSON("msg")
	out := fmt("something happened")
	if !strings.Contains(out, `"msg"`) {
		t.Fatalf("expected msg key in output, got %q", out)
	}
	if !strings.Contains(out, `"something happened"`) {
		t.Fatalf("expected line value in output, got %q", out)
	}
}

func TestJSON_AlreadyJSON(t *testing.T) {
	fmt := format.JSON("msg")
	raw := `{"level":"info","text":"ok"}`
	out := fmt(raw)
	if out != raw {
		t.Fatalf("expected passthrough of valid JSON, got %q", out)
	}
}

func TestJSON_DefaultKey(t *testing.T) {
	fmt := format.JSON("")
	out := fmt("test")
	if !strings.Contains(out, `"msg"`) {
		t.Fatalf("expected default key 'msg', got %q", out)
	}
}

func TestLogfmt_NoSpaces(t *testing.T) {
	fmt := format.Logfmt("msg")
	out := fmt("connected")
	if out != "msg=connected" {
		t.Fatalf("expected 'msg=connected', got %q", out)
	}
}

func TestLogfmt_WithSpaces(t *testing.T) {
	fmt := format.Logfmt("msg")
	out := fmt("hello world")
	if !strings.HasPrefix(out, "msg=") {
		t.Fatalf("expected logfmt key prefix, got %q", out)
	}
	if !strings.Contains(out, `"hello world"`) {
		t.Fatalf("expected quoted value, got %q", out)
	}
}

func TestTemplate_Substitution(t *testing.T) {
	fmt := format.Template("[LOG] {line}")
	out := fmt("disk full")
	if out != "[LOG] disk full" {
		t.Fatalf("expected '[LOG] disk full', got %q", out)
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	f := format.Chain(
		format.Template("prefix:{line}"),
		format.Template("{line}:suffix"),
	)
	out := f("data")
	if out != "prefix:data:suffix" {
		t.Fatalf("expected 'prefix:data:suffix', got %q", out)
	}
}

func TestNew_ReturnsCorrectFormatter(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		input   string
		contains string
	}{
		{"json", "msg", "hi", `"msg"`},
		{"logfmt", "msg", "hi", "msg=hi"},
		{"passthrough", "", "hi", "hi"},
		{"unknown", "", "hi", "hi"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := format.New(tc.name, tc.key)
			out := f(tc.input)
			if !strings.Contains(out, tc.contains) {
				t.Fatalf("name=%s: expected output to contain %q, got %q", tc.name, tc.contains, out)
			}
		})
	}
}
