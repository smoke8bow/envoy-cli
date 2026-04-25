package envflag_test

import (
	"errors"
	"testing"

	"github.com/envoy-cli/envoy/internal/envflag"
)

// fakeStore is an in-memory Store implementation.
type fakeStore struct {
	data map[string]map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]map[string]string)}
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	f.data[name] = vars
	return nil
}

func newManager(t *testing.T) *envflag.Manager {
	t.Helper()
	m, err := envflag.NewManager(newFakeStore(), "flags")
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestNewManagerEmptyProfile(t *testing.T) {
	_, err := envflag.NewManager(newFakeStore(), "")
	if err == nil {
		t.Fatal("expected error for empty profile")
	}
}

func TestSetAndIsEnabled(t *testing.T) {
	m := newManager(t)
	if err := m.Set("DARK_MODE", true); err != nil {
		t.Fatalf("Set: %v", err)
	}
	ok, err := m.IsEnabled("DARK_MODE")
	if err != nil {
		t.Fatalf("IsEnabled: %v", err)
	}
	if !ok {
		t.Error("expected DARK_MODE to be enabled")
	}
}

func TestSetDisable(t *testing.T) {
	m := newManager(t)
	_ = m.Set("BETA", true)
	_ = m.Set("BETA", false)
	ok, err := m.IsEnabled("BETA")
	if err != nil {
		t.Fatalf("IsEnabled: %v", err)
	}
	if ok {
		t.Error("expected BETA to be disabled")
	}
}

func TestIsEnabledMissingFlag(t *testing.T) {
	m := newManager(t)
	ok, err := m.IsEnabled("MISSING")
	if err != nil {
		t.Fatalf("IsEnabled: %v", err)
	}
	if ok {
		t.Error("expected false for missing flag")
	}
}

func TestDelete(t *testing.T) {
	m := newManager(t)
	_ = m.Set("FEATURE_X", true)
	if err := m.Delete("FEATURE_X"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	ok, _ := m.IsEnabled("FEATURE_X")
	if ok {
		t.Error("expected flag to be gone after delete")
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newManager(t)
	if err := m.Delete("GHOST"); err == nil {
		t.Fatal("expected error deleting non-existent flag")
	}
}

func TestList(t *testing.T) {
	m := newManager(t)
	_ = m.Set("A", true)
	_ = m.Set("B", false)
	flags, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(flags))
	}
	if !flags["A"] {
		t.Error("expected A=true")
	}
	if flags["B"] {
		t.Error("expected B=false")
	}
}

func TestSetEmptyFlagName(t *testing.T) {
	m := newManager(t)
	if err := m.Set("", true); err == nil {
		t.Fatal("expected error for empty flag name")
	}
}
