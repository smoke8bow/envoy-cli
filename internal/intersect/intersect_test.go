package intersect_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/intersect"
)

func TestIntersectNoCommonKeys(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3", "QUX": "4"}
	r := intersect.Intersect(a, b)
	if len(r.Vars) != 0 {
		t.Fatalf("expected 0 vars, got %d", len(r.Vars))
	}
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected 0 conflicts, got %d", len(r.Conflicts))
	}
}

func TestIntersectCommonKeysNoConflict(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2", "ONLY_A": "x"}
	b := map[string]string{"FOO": "1", "BAR": "2", "ONLY_B": "y"}
	r := intersect.Intersect(a, b)
	if len(r.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(r.Vars))
	}
	if r.Vars["FOO"] != "1" || r.Vars["BAR"] != "2" {
		t.Fatalf("unexpected vars: %v", r.Vars)
	}
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %v", r.Conflicts)
	}
}

func TestIntersectConflicts(t *testing.T) {
	a := map[string]string{"FOO": "alpha", "BAR": "same"}
	b := map[string]string{"FOO": "beta", "BAR": "same"}
	r := intersect.Intersect(a, b)
	if len(r.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(r.Vars))
	}
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	c, ok := r.Conflicts["FOO"]
	if !ok {
		t.Fatal("expected conflict on FOO")
	}
	if c.ValueA != "alpha" || c.ValueB != "beta" {
		t.Fatalf("unexpected conflict values: %+v", c)
	}
}

func TestIntersectValuesFromA(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	r := intersect.Intersect(a, b)
	if r.Vars["KEY"] != "from-a" {
		t.Fatalf("expected value from a, got %q", r.Vars["KEY"])
	}
}

func TestKeysSorted(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "2", "M": "3"}
	b := map[string]string{"Z": "1", "A": "2", "M": "3"}
	r := intersect.Intersect(a, b)
	keys := r.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("keys not sorted: got %v", keys)
		}
	}
}

func TestConflictKeysSorted(t *testing.T) {
	a := map[string]string{"Z": "z1", "A": "a1", "M": "same"}
	b := map[string]string{"Z": "z2", "A": "a2", "M": "same"}
	r := intersect.Intersect(a, b)
	keys := r.ConflictKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 conflict keys, got %v", keys)
	}
	if keys[0] != "A" || keys[1] != "Z" {
		t.Fatalf("conflict keys not sorted: %v", keys)
	}
}
