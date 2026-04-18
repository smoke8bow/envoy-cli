package trim

import (
	"strings"
	"testing"
)

func newTrimmer(max int) *Trimmer { return NewTrimmer(max) }

func TestApplyNoOversized(t *testing.T) {
	tr := newTrimmer(10)
	vars := map[string]string{"A": "hello", "B": "world"}
	r := tr.Apply(vars)
	if len(r.Removed) != 0 {
		t.Fatalf("expected no removed, got %v", r.Removed)
	}
	if len(r.Kept) != 2 {
		t.Fatalf("expected 2 kept, got %d", len(r.Kept))
	}
}

func TestApplyRemovesOversized(t *testing.T) {
	tr := newTrimmer(5)
	vars := map[string]string{"SHORT": "hi", "LONG": "this is too long"}
	r := tr.Apply(vars)
	if len(r.Removed) != 1 || r.Removed[0] != "LONG" {
		t.Fatalf("expected LONG removed, got %v", r.Removed)
	}
	if _, ok := r.Kept["SHORT"]; !ok {
		t.Fatal("expected SHORT in kept")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tr := newTrimmer(3)
	vars := map[string]string{"K": "abcdef"}
	tr.Apply(vars)
	if vars["K"] != "abcdef" {
		t.Fatal("input was mutated")
	}
}

func TestTrimValues(t *testing.T) {
	tr := newTrimmer(4)
	vars := map[string]string{"A": "hello", "B": "hi"}
	out := tr.TrimValues(vars)
	if out["A"] != "hell" {
		t.Fatalf("expected truncated value, got %q", out["A"])
	}
	if out["B"] != "hi" {
		t.Fatalf("expected unchanged value, got %q", out["B"])
	}
}

func TestDefaultMaxLen(t *testing.T) {
	tr := NewTrimmer(0)
	if tr.maxLen != 256 {
		t.Fatalf("expected default 256, got %d", tr.maxLen)
	}
}

func TestFormatNoRemoved(t *testing.T) {
	r := Result{Removed: nil, Kept: map[string]string{}}
	if Format(r) != "no oversized variables found" {
		t.Fatal("unexpected format output")
	}
}

func TestFormatWithRemoved(t *testing.T) {
	r := Result{Removed: []string{"FOO", "BAR"}}
	out := Format(r)
	if !strings.Contains(out, "2") || !strings.Contains(out, "FOO") {
		t.Fatalf("unexpected format output: %s", out)
	}
}
