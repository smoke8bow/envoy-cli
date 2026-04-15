package profile_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/profile"
	"envoy-cli/internal/store"
)

func newManager(t *testing.T) *profile.Manager {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.json")
	s, err := store.Load(path)
	if err != nil {
		t.Fatalf("store.Load: %v", err)
	}
	return profile.NewManager(s)
}

func TestCreateAndGet(t *testing.T) {
	m := newManager(t)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}

	if err := m.Create("dev", vars); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := m.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", got)
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := newManager(t)
	m.Create("dev", map[string]string{"X": "1"})
	err := m.Create("dev", map[string]string{"X": "2"})
	if !errors.Is(err, profile.ErrProfileExists) {
		t.Errorf("expected ErrProfileExists, got %v", err)
	}
}

func TestGetNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Get("missing")
	if !errors.Is(err, profile.ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	m := newManager(t)
	m.Create("prod", map[string]string{"ENV": "production"})

	if err := m.Update("prod", map[string]string{"ENV": "staging"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := m.Get("prod")
	if got["ENV"] != "staging" {
		t.Errorf("expected staging, got %s", got["ENV"])
	}
}

func TestUpdateNotFound(t *testing.T) {
	m := newManager(t)
	err := m.Update("ghost", map[string]string{})
	if !errors.Is(err, profile.ErrProfileNotFound) {
		t.Errorf("expected ErrProfileNotFound, got %v", err)
	}
}

func TestDeleteAndList(t *testing.T) {
	m := newManager(t)
	m.Create("a", map[string]string{})
	m.Create("b", map[string]string{})

	if err := m.Delete("a"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	names := m.List()
	if len(names) != 1 || names[0] != "b" {
		t.Errorf("expected [b], got %v", names)
	}
	_ = os.Getenv // suppress unused import lint
}
