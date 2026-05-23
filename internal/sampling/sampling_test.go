package sampling_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/sampling"
)

func TestNew_InvalidRate(t *testing.T) {
	cases := []float64{0, -0.5, 1.1, 2}
	for _, r := range cases {
		_, err := sampling.New(sampling.Config{Rate: r})
		if err == nil {
			t.Errorf("expected error for rate %.2f, got nil", r)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	_, err := sampling.New(sampling.Config{Rate: 0.5, Seed: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSampler_KeepsAll(t *testing.T) {
	s, err := sampling.New(sampling.Config{Rate: 1.0, Seed: 1})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if !s.Sample("some log line") {
			t.Fatal("rate=1.0 should keep every line")
		}
	}
}

func TestSampler_ApproximateRate(t *testing.T) {
	const n = 100_000
	s, err := sampling.New(sampling.Config{Rate: 0.1, Seed: 99})
	if err != nil {
		t.Fatal(err)
	}
	kept := 0
	for i := 0; i < n; i++ {
		if s.Sample("line") {
			kept++
		}
	}
	got := float64(kept) / n
	if got < 0.08 || got > 0.12 {
		t.Errorf("expected ~10%% kept, got %.2f%%", got*100)
	}
}

func TestPassThrough_KeepsAll(t *testing.T) {
	s := sampling.PassThrough()
	for i := 0; i < 50; i++ {
		if !s.Sample("x") {
			t.Fatal("PassThrough should keep all lines")
		}
	}
}

func TestSampler_Stats(t *testing.T) {
	type statter interface {
		Stats() (int64, int64)
	}
	s, _ := sampling.New(sampling.Config{Rate: 1.0, Seed: 0})
	st, ok := s.(statter)
	if !ok {
		t.Skip("sampler does not expose Stats")
	}
	for i := 0; i < 10; i++ {
		s.Sample("line")
	}
	total, sampled := st.Stats()
	if total != 10 {
		t.Errorf("total: want 10, got %d", total)
	}
	if sampled != 10 {
		t.Errorf("sampled: want 10, got %d", sampled)
	}
}
