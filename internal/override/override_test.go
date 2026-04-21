package override_test

import (
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/override"
)

func newManager() *override.Manager {
	return override.NewManager()
}

func TestSetAndApply(t *testing.T) {
	m := newManager()
	base := map[string]string{"FOO": "foo", "BAR": "bar"}
	if err := m.Set("dev", "FOO", "overridden"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := m.Apply("dev", base)
	if result["FOO"] != "overridden" {
		t.Errorf("expected overridden, got %q", result["FOO"])
	}
	if result["BAR"] != "bar" {
		t.Errorf("expected bar, got %q", result["BAR"])
	}
}

func TestApplyDoesNotMutateBase(t *testing.T) {
	m := newManager()
	base := map[string]string{"KEY": "original"}
	_ = m.Set("dev", "KEY", "changed")
	_ = m.Apply("dev", base)
	if base["KEY"] != "original" {
		t.Error("Apply mutated the base map")
	}
}

func TestUnset(t *testing.T) {
	m := newManager()
	_ = m.Set("dev", "FOO", "val")
	if err := m.Unset("dev", "FOO"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	layer := m.Layer("dev")
	if _, ok := layer["FOO"]; ok {
		t.Error("expected FOO to be removed from layer")
	}
}

func TestUnsetNotFound(t *testing.T) {
	m := newManager()
	err := m.Unset("dev", "MISSING")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestSetEmptyProfileError(t *testing.T) {
	m := newManager()
	if err := m.Set("", "KEY", "val"); err == nil {
		t.Error("expected error for empty profile")
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	m := newManager()
	if err := m.Set("dev", "", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestClear(t *testing.T) {
	m := newManager()
	_ = m.Set("dev", "A", "1")
	_ = m.Set("dev", "B", "2")
	m.Clear("dev")
	layer := m.Layer("dev")
	if len(layer) != 0 {
		t.Errorf("expected empty layer after Clear, got %d entries", len(layer))
	}
}

func TestApplyNoLayerReturnsBaseClone(t *testing.T) {
	m := newManager()
	base := map[string]string{"X": "1"}
	result := m.Apply("nonexistent", base)
	if result["X"] != "1" {
		t.Errorf("expected 1, got %q", result["X"])
	}
}

func TestLayerReturnsCopy(t *testing.T) {
	m := newManager()
	_ = m.Set("dev", "FOO", "bar")
	layer := m.Layer("dev")
	layer["FOO"] = "mutated"
	original := m.Layer("dev")
	if original["FOO"] != "bar" {
		t.Error("Layer did not return a copy; internal state was mutated")
	}
}
