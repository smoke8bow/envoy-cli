package prefix

import (
	"testing"
)

func newManager() *Manager { return NewManager() }

func TestApplyAddsPrefix(t *testing.T) {
	m := newManager()
	out, err := m.Apply(map[string]string{"FOO": "bar", "BAZ": "qux"}, "APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_FOO"] != "bar" || out["APP_BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestApplyEmptyPrefixError(t *testing.T) {
	m := newManager()
	_, err := m.Apply(map[string]string{"K": "v"}, "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestStripRemovesPrefix(t *testing.T) {
	m := newManager()
	out, err := m.Strip(map[string]string{"APP_FOO": "bar", "APP_BAZ": "qux"}, "APP_", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestStripKeepsUnprefixedWhenFlagFalse(t *testing.T) {
	m := newManager()
	out, _ := m.Strip(map[string]string{"APP_FOO": "bar", "OTHER": "val"}, "APP_", false)
	if _, ok := out["OTHER"]; !ok {
		t.Fatal("expected OTHER to be kept")
	}
}

func TestStripOmitsUnprefixedWhenFlagTrue(t *testing.T) {
	m := newManager()
	out, _ := m.Strip(map[string]string{"APP_FOO": "bar", "OTHER": "val"}, "APP_", true)
	if _, ok := out["OTHER"]; ok {
		t.Fatal("expected OTHER to be omitted")
	}
}

func TestFilterReturnsOnlyMatching(t *testing.T) {
	m := newManager()
	out := m.Filter(map[string]string{"APP_A": "1", "DB_B": "2", "APP_C": "3"}, "APP_")
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DB_B"]; ok {
		t.Fatal("DB_B should not be in result")
	}
}

func TestFilterEmptyResult(t *testing.T) {
	m := newManager()
	out := m.Filter(map[string]string{"FOO": "bar"}, "NOPE_")
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}
