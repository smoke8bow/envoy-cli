package truncate

import (
	"strings"
	"testing"
)

func newTruncator(maxLen int) *Truncator {
	return NewTruncator(Options{MaxLen: maxLen, Suffix: "..."})
}

func TestApplyNoTruncation(t *testing.T) {
	tr := newTruncator(20)
	vars := map[string]string{"KEY": "short"}
	out, results := tr.Apply(vars)
	if out["KEY"] != "short" {
		t.Fatalf("expected 'short', got %q", out["KEY"])
	}
	if results[0].WasCut {
		t.Fatal("expected WasCut=false")
	}
}

func TestApplyTruncates(t *testing.T) {
	tr := newTruncator(10)
	vars := map[string]string{"KEY": "this is a very long value"}
	out, results := tr.Apply(vars)
	if len(out["KEY"]) > 10 {
		t.Fatalf("value too long: %q", out["KEY"])
	}
	if !strings.HasSuffix(out["KEY"], "...") {
		t.Fatalf("expected suffix '...', got %q", out["KEY"])
	}
	if !results[0].WasCut {
		t.Fatal("expected WasCut=true")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tr := newTruncator(5)
	vars := map[string]string{"A": "hello world"}
	original := vars["A"]
	tr.Apply(vars)
	if vars["A"] != original {
		t.Fatal("input map was mutated")
	}
}

func TestKeepAll(t *testing.T) {
	tr := NewTruncator(Options{MaxLen: 5, Suffix: "...", KeepAll: true})
	vars := map[string]string{"K": "this should not be truncated"}
	out, results := tr.Apply(vars)
	if out["K"] != vars["K"] {
		t.Fatal("KeepAll should prevent truncation")
	}
	if results[0].WasCut {
		t.Fatal("expected WasCut=false with KeepAll")
	}
}

func TestFormat(t *testing.T) {
	results := []Result{
		{WasCut: true},
		{WasCut: false},
		{WasCut: true},
	}
	out := Format(results)
	if !strings.Contains(out, "2") {
		t.Fatalf("expected '2' in format output, got %q", out)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxLen != 64 {
		t.Fatalf("expected MaxLen=64, got %d", opts.MaxLen)
	}
}
