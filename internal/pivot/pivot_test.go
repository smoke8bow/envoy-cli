package pivot

import (
	"testing"
)

func TestIsSupportedValid(t *testing.T) {
	for _, d := range Supported() {
		if !IsSupported(d) {
			t.Errorf("expected %q to be supported", d)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if IsSupported(Direction("unknown")) {
		t.Error("expected unknown direction to be unsupported")
	}
}

func TestPivotKeysToValues(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Pivot(vars, DirectionKeysToValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["bar"] != "FOO" {
		t.Errorf("expected out[bar]=FOO, got %q", out["bar"])
	}
	if out["qux"] != "BAZ" {
		t.Errorf("expected out[qux]=BAZ, got %q", out["qux"])
	}
}

func TestPivotValuesToKeys(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	out, err := Pivot(vars, DirectionValuesToKeys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["1"] != "A" {
		t.Errorf("expected out[1]=A, got %q", out["1"])
	}
}

func TestPivotDuplicateValues(t *testing.T) {
	vars := map[string]string{"X": "same", "Y": "same"}
	_, err := Pivot(vars, DirectionKeysToValues)
	if err == nil {
		t.Error("expected error for duplicate values")
	}
}

func TestPivotUnsupportedDirection(t *testing.T) {
	_, err := Pivot(map[string]string{"K": "V"}, Direction("sideways"))
	if err == nil {
		t.Error("expected error for unsupported direction")
	}
}

func TestPivotEmpty(t *testing.T) {
	out, err := Pivot(map[string]string{}, DirectionKeysToValues)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
