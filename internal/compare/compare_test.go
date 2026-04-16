package compare

import (
	"strings"
	"testing"
)

func TestCompareIdentical(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(a, b)
	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 || len(r.Different) != 0 {
		t.Fatalf("expected no differences, got %+v", r)
	}
	if len(r.Same) != 2 {
		t.Fatalf("expected 2 same keys, got %d", len(r.Same))
	}
}

func TestCompareOnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "x"}
	b := map[string]string{"FOO": "bar"}
	r := Compare(a, b)
	if _, ok := r.OnlyInA["ONLY_A"]; !ok {
		t.Fatal("expected ONLY_A in OnlyInA")
	}
}

func TestCompareOnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "y"}
	r := Compare(a, b)
	if _, ok := r.OnlyInB["ONLY_B"]; !ok {
		t.Fatal("expected ONLY_B in OnlyInB")
	}
}

func TestCompareDifferent(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Compare(a, b)
	p, ok := r.Different["FOO"]
	if !ok {
		t.Fatal("expected FOO in Different")
	}
	if p.A != "old" || p.B != "new" {
		t.Fatalf("unexpected pair %+v", p)
	}
}

func TestFormat(t *testing.T) {
	a := map[string]string{"FOO": "old", "GONE": "x"}
	b := map[string]string{"FOO": "new", "NEW": "y"}
	r := Compare(a, b)
	out := Format("staging", "prod", r)
	if !strings.Contains(out, "staging") || !strings.Contains(out, "prod") {
		t.Fatal("format missing profile names")
	}
	if !strings.Contains(out, "~ FOO") {
		t.Fatal("format missing changed key")
	}
	if !strings.Contains(out, "< GONE") {
		t.Fatal("format missing only-in-A key")
	}
	if !strings.Contains(out, "> NEW") {
		t.Fatal("format missing only-in-B key")
	}
}
