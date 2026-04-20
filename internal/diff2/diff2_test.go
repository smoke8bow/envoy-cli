package diff2_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/diff2"
)

func TestComputeOnlyLeft(t *testing.T) {
	r := diff2.Compute(map[string]string{"A": "1"}, map[string]string{})
	if len(r.OnlyLeft()) != 1 || r.OnlyLeft()[0].Key != "A" {
		t.Fatal("expected A in left-only")
	}
}

func TestComputeOnlyRight(t *testing.T) {
	r := diff2.Compute(map[string]string{}, map[string]string{"B": "2"})
	if len(r.OnlyRight()) != 1 || r.OnlyRight()[0].Key != "B" {
		t.Fatal("expected B in right-only")
	}
}

func TestComputeChanged(t *testing.T) {
	r := diff2.Compute(map[string]string{"X": "old"}, map[string]string{"X": "new"})
	if len(r.Changed()) != 1 {
		t.Fatal("expected one changed entry")
	}
	e := r.Changed()[0]
	if e.Left != "old" || e.Right != "new" {
		t.Fatalf("unexpected values: %+v", e)
	}
}

func TestComputeUnchanged(t *testing.T) {
	r := diff2.Compute(map[string]string{"K": "v"}, map[string]string{"K": "v"})
	if len(r.Unchanged()) != 1 {
		t.Fatal("expected one unchanged entry")
	}
	if !r.Unchanged()[0].Equal {
		t.Fatal("entry should be equal")
	}
}

func TestComputeSorted(t *testing.T) {
	r := diff2.Compute(
		map[string]string{"Z": "1", "A": "2", "M": "3"},
		map[string]string{"Z": "1", "A": "2", "M": "3"},
	)
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Fatalf("expected sorted keys, got %v", keys)
	}
}

func TestFormatOnlyLeft(t *testing.T) {
	r := diff2.Compute(map[string]string{"FOO": "bar"}, map[string]string{})
	out := diff2.Format(r, diff2.DefaultFormatOptions())
	if !strings.Contains(out, "- [left only] FOO=bar") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFormatChanged(t *testing.T) {
	r := diff2.Compute(map[string]string{"KEY": "old"}, map[string]string{"KEY": "new"})
	opts := diff2.FormatOptions{LeftLabel: "prod", RightLabel: "staging"}
	out := diff2.Format(r, opts)
	if !strings.Contains(out, "~ KEY: prod=old → staging=new") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestSummaryNoDiff(t *testing.T) {
	r := diff2.Compute(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	if diff2.Summary(r) != "no differences" {
		t.Fatal("expected no differences")
	}
}

func TestSummaryMixed(t *testing.T) {
	r := diff2.Compute(
		map[string]string{"A": "1", "B": "old"},
		map[string]string{"B": "new", "C": "3"},
	)
	s := diff2.Summary(r)
	if !strings.Contains(s, "added") || !strings.Contains(s, "removed") || !strings.Contains(s, "changed") {
		t.Fatalf("unexpected summary: %q", s)
	}
}
