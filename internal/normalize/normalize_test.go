package normalize_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/normalize"
)

func TestPassThrough(t *testing.T) {
	line := "  hello world  "
	if got := normalize.PassThrough(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestTrimSpace(t *testing.T) {
	cases := []struct{ in, want string }{
		{"  hello  ", "hello"},
		{"\t\nfoo\n", "foo"},
		{"no-trim", "no-trim"},
		{"", ""},
	}
	for _, c := range cases {
		if got := normalize.TrimSpace(c.in); got != c.want {
			t.Errorf("TrimSpace(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCollapseWhitespace(t *testing.T) {
	cases := []struct{ in, want string }{
		{"foo  bar", "foo bar"},
		{"a\t\tb", "a b"},
		{"single", "single"},
		{"  leading", " leading"},
		{"", ""},
	}
	for _, c := range cases {
		if got := normalize.CollapseWhitespace(c.in); got != c.want {
			t.Errorf("CollapseWhitespace(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestStripControl(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello\x00world", "helloworld"},
		{"tab\there", "tab\there"},
		{"bell\x07end", "bellend"},
		{"clean", "clean"},
	}
	for _, c := range cases {
		if got := normalize.StripControl(c.in); got != c.want {
			t.Errorf("StripControl(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeNewlines(t *testing.T) {
	cases := []struct{ in, want string }{
		{"foo\r\nbar", "foo\nbar"},
		{"foo\rbar", "foo\nbar"},
		{"foo\nbar", "foo\nbar"},
		{"no newline", "no newline"},
	}
	for _, c := range cases {
		if got := normalize.NormalizeNewlines(c.in); got != c.want {
			t.Errorf("NormalizeNewlines(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	fn := normalize.Chain(
		normalize.NormalizeNewlines,
		normalize.TrimSpace,
		normalize.CollapseWhitespace,
	)
	input := "  foo\r\n  bar  baz  "
	want := "foo\n  bar  baz"
	if got := fn(input); got != want {
		t.Errorf("Chain result = %q, want %q", got, want)
	}
}
