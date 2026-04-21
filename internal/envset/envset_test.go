package envset_test

import (
	"testing"

	"envoy-cli/internal/envset"
)

func TestIsSupportedValid(t *testing.T) {
	for _, op := range envset.Supported() {
		if !envset.IsSupported(op) {
			t.Errorf("expected %q to be supported", op)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if envset.IsSupported("bogus") {
		t.Error("expected bogus to be unsupported")
	}
}

func TestUnion(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "99", "C": "3"}
	got := envset.Union(a, b)
	if got["A"] != "1" || got["B"] != "99" || got["C"] != "3" {
		t.Errorf("unexpected union result: %v", got)
	}
}

func TestUnionDoesNotMutate(t *testing.T) {
	a := map[string]string{"X": "orig"}
	b := map[string]string{"X": "new"}
	_ = envset.Union(a, b)
	if a["X"] != "orig" {
		t.Error("Union mutated input a")
	}
}

func TestIntersection(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "99", "C": "3"}
	got := envset.Intersection(a, b)
	if len(got) != 1 || got["B"] != "99" {
		t.Errorf("unexpected intersection result: %v", got)
	}
}

func TestIntersectionEmpty(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	got := envset.Intersection(a, b)
	if len(got) != 0 {
		t.Errorf("expected empty intersection, got %v", got)
	}
}

func TestDifference(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"B": "2"}
	got := envset.Difference(a, b)
	if len(got) != 2 || got["A"] != "1" || got["C"] != "3" {
		t.Errorf("unexpected difference result: %v", got)
	}
}

func TestApplyUnknownOp(t *testing.T) {
	_, err := envset.Apply("nope", nil, nil)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApplyDispatch(t *testing.T) {
	a := map[string]string{"K": "v"}
	b := map[string]string{"K": "v2", "L": "x"}

	for _, tc := range []struct {
		op      envset.Op
		wantLen int
	}{
		{envset.OpUnion, 2},
		{envset.OpIntersection, 1},
		{envset.OpDifference, 0},
	} {
		got, err := envset.Apply(tc.op, a, b)
		if err != nil {
			t.Fatalf("op %s: unexpected error: %v", tc.op, err)
		}
		if len(got) != tc.wantLen {
			t.Errorf("op %s: want len %d, got %d", tc.op, tc.wantLen, len(got))
		}
	}
}
