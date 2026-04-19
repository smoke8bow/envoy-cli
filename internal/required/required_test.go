package required

import (
	"testing"
)

func newManager() *Manager {
	return NewManager()
}

func TestSetAndGet(t *testing.T) {
	m := newManager()
	m.Set("dev", []string{"FOO", "BAR"})
	keys := m.Get("dev")
	if len(keys) != 2 || keys[0] != "FOO" || keys[1] != "BAR" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestCheckNoViolations(t *testing.T) {
	m := newManager()
	m.Set("prod", []string{"DB_URL", "SECRET"})
	vars := map[string]string{"DB_URL": "postgres://", "SECRET": "abc"}
	if v := m.Check("prod", vars); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestCheckMissingKey(t *testing.T) {
	m := newManager()
	m.Set("prod", []string{"DB_URL", "SECRET"})
	vars := map[string]string{"DB_URL": "postgres://"}
	v := m.Check("prod", vars)
	if len(v) != 1 || v[0].Key != "SECRET" {
		t.Fatalf("expected SECRET violation, got %v", v)
	}
}

func TestCheckEmptyValue(t *testing.T) {
	m := newManager()
	m.Set("prod", []string{"API_KEY"})
	vars := map[string]string{"API_KEY": ""}
	v := m.Check("prod", vars)
	if len(v) != 1 || v[0].Key != "API_KEY" {
		t.Fatalf("expected API_KEY violation, got %v", v)
	}
}

func TestCheckUnknownProfile(t *testing.T) {
	m := newManager()
	v := m.Check("ghost", map[string]string{"X": "y"})
	if len(v) != 0 {
		t.Fatalf("expected no violations for unknown profile")
	}
}

func TestViolationError(t *testing.T) {
	v := Violation{Key: "FOO", Profile: "prod"}
	if v.Error() == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestCheckAll(t *testing.T) {
	m := newManager()
	m.Set("a", []string{"K1"})
	m.Set("b", []string{"K2"})
	store := map[string]map[string]string{
		"a": {"K1": "val"},
		"b": {},
	}
	result := m.CheckAll(func(p string) map[string]string { return store[p] })
	if _, bad := result["b"]; !bad {
		t.Fatal("expected violation for profile b")
	}
	if _, good := result["a"]; good {
		t.Fatal("expected no violation for profile a")
	}
}
