package truncate

import (
	"strings"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	tr, err := New(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.maxBytes != DefaultMaxBytes {
		t.Errorf("maxBytes: got %d, want %d", tr.maxBytes, DefaultMaxBytes)
	}
	if tr.suffix != DefaultSuffix {
		t.Errorf("suffix: got %q, want %q", tr.suffix, DefaultSuffix)
	}
}

func TestNew_NegativeMaxBytes(t *testing.T) {
	_, err := New(Config{MaxBytes: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxBytes, got nil")
	}
}

func TestTransform_ShortLine(t *testing.T) {
	tr, _ := New(Config{MaxBytes: 20, Suffix: "..."})
	input := "hello world"
	got := tr.Transform(input)
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestTransform_ExactLength(t *testing.T) {
	tr, _ := New(Config{MaxBytes: 5, Suffix: "..."})
	input := "hello"
	got := tr.Transform(input)
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestTransform_TruncatesLongLine(t *testing.T) {
	tr, _ := New(Config{MaxBytes: 10, Suffix: "..."})
	input := "hello world extra"
	got := tr.Transform(input)
	if len(got) != 10 {
		t.Errorf("length: got %d, want 10", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestTransform_SuffixLongerThanMax(t *testing.T) {
	// suffix is longer than MaxBytes; result should be trimmed suffix
	tr, _ := New(Config{MaxBytes: 3, Suffix: "[truncated]"})
	input := "some very long line that should be cut"
	got := tr.Transform(input)
	if len(got) > 3 {
		t.Errorf("expected len <= 3, got %d (%q)", len(got), got)
	}
}

func TestTransform_CustomSuffix(t *testing.T) {
	tr, _ := New(Config{MaxBytes: 15, Suffix: "[CUT]"})
	input := "abcdefghijklmnopqrstuvwxyz"
	got := tr.Transform(input)
	if !strings.HasSuffix(got, "[CUT]") {
		t.Errorf("expected [CUT] suffix, got %q", got)
	}
	if len(got) != 15 {
		t.Errorf("expected length 15, got %d", len(got))
	}
}

func TestAsTransformFunc(t *testing.T) {
	tr, _ := New(Config{MaxBytes: 8, Suffix: "..."})
	fn := tr.AsTransformFunc()
	got := fn("hello world")
	if len(got) != 8 {
		t.Errorf("transform func: got len %d, want 8", len(got))
	}
}
