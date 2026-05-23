package routing_test

import (
	"sort"
	"testing"

	"github.com/yourorg/logpipe/internal/routing"
)

func TestNew_ValidRoutes(t *testing.T) {
	r, err := routing.New([]routing.Route{
		{Source: "app", Sinks: []string{"stdout", "file"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks, ok := r.Resolve("app")
	if !ok {
		t.Fatal("expected route for 'app'")
	}
	if len(sinks) != 2 {
		t.Fatalf("expected 2 sinks, got %d", len(sinks))
	}
}

func TestNew_EmptySource(t *testing.T) {
	_, err := routing.New([]routing.Route{
		{Source: "", Sinks: []string{"stdout"}},
	})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNew_NoSinks(t *testing.T) {
	_, err := routing.New([]routing.Route{
		{Source: "app", Sinks: nil},
	})
	if err == nil {
		t.Fatal("expected error for missing sinks")
	}
}

func TestResolve_UnknownSource(t *testing.T) {
	r, _ := routing.New([]routing.Route{
		{Source: "app", Sinks: []string{"stdout"}},
	})
	_, ok := r.Resolve("unknown")
	if ok {
		t.Fatal("expected no route for unknown source")
	}
}

func TestAdd_AppendsRoute(t *testing.T) {
	r, _ := routing.New(nil)
	if err := r.Add("svc", []string{"file"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sinks, ok := r.Resolve("svc")
	if !ok || len(sinks) != 1 || sinks[0] != "file" {
		t.Fatalf("unexpected sinks: %v", sinks)
	}
}

func TestAdd_EmptySource(t *testing.T) {
	r, _ := routing.New(nil)
	if err := r.Add("", []string{"file"}); err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestSources_ReturnsAll(t *testing.T) {
	r, _ := routing.New([]routing.Route{
		{Source: "a", Sinks: []string{"s1"}},
		{Source: "b", Sinks: []string{"s2"}},
	})
	srcs := r.Sources()
	sort.Strings(srcs)
	if len(srcs) != 2 || srcs[0] != "a" || srcs[1] != "b" {
		t.Fatalf("unexpected sources: %v", srcs)
	}
}

func TestResolve_ReturnsCopy(t *testing.T) {
	r, _ := routing.New([]routing.Route{
		{Source: "app", Sinks: []string{"stdout"}},
	})
	sinks, _ := r.Resolve("app")
	sinks[0] = "mutated"
	original, _ := r.Resolve("app")
	if original[0] == "mutated" {
		t.Fatal("Resolve should return a copy, not a reference")
	}
}
