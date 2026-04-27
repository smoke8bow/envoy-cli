package cascade

import (
	"errors"
	"testing"
)

// fakeStore is an in-memory Accessor for testing.
type fakeStore struct {
	profiles map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func newManager(profiles map[string]map[string]string) *Manager {
	return NewManager(&fakeStore{profiles: profiles})
}

func TestResolveSingleProfile(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"A": "1", "B": "2"},
	})
	res, err := m.Resolve([]string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
}

func TestResolveLaterOverrides(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base":     {"A": "base", "B": "base"},
		"override": {"B": "override", "C": "override"},
	})
	res, err := m.Resolve([]string{"base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "base" {
		t.Errorf("A should come from base, got %q", res.Vars["A"])
	}
	if res.Vars["B"] != "override" {
		t.Errorf("B should be overridden, got %q", res.Vars["B"])
	}
	if res.Source["B"] != "override" {
		t.Errorf("source for B should be 'override', got %q", res.Source["B"])
	}
}

func TestResolveSourceTracking(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"X": "1"},
		"mid":  {"Y": "2"},
		"top":  {"X": "3"},
	})
	res, err := m.Resolve([]string{"base", "mid", "top"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source["X"] != "top" {
		t.Errorf("expected source 'top' for X, got %q", res.Source["X"])
	}
	if res.Source["Y"] != "mid" {
		t.Errorf("expected source 'mid' for Y, got %q", res.Source["Y"])
	}
}

func TestResolveEmptyProfilesError(t *testing.T) {
	m := newManager(nil)
	_, err := m.Resolve([]string{})
	if err == nil {
		t.Fatal("expected error for empty profiles")
	}
}

func TestResolveUnknownProfileError(t *testing.T) {
	m := newManager(map[string]map[string]string{})
	_, err := m.Resolve([]string{"missing"})
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestResultKeys(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"p": {"Z": "z", "A": "a", "M": "m"},
	})
	res, _ := m.Resolve([]string{"p"})
	keys := res.Keys()
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

// TestResolveSingleProfileSourceTracking verifies that source tracking is
// correct when only a single profile is resolved — every key should point
// back to that profile.
func TestResolveSingleProfileSourceTracking(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"A": "1", "B": "2"},
	})
	res, err := m.Resolve([]string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"A", "B"} {
		if res.Source[key] != "base" {
			t.Errorf("expected source 'base' for key %q, got %q", key, res.Source[key])
		}
	}
}
