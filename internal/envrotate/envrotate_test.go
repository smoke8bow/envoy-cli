package envrotate_test

import (
	"errors"
	"testing"

	"github.com/your-org/envoy-cli/internal/envrotate"
)

// fakeStore is an in-memory Store for testing.
type fakeStore struct {
	profiles map[string]map[string]string
	setErr   error
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.profiles[name] = vars
	return nil
}

func newManager(t *testing.T, profiles map[string]map[string]string) (*envrotate.Manager, *fakeStore) {
	t.Helper()
	store := &fakeStore{profiles: profiles}
	return envrotate.NewManager(store, envrotate.DefaultOptions()), store
}

func TestRotateRenamesKeys(t *testing.T) {
	m, store := newManager(t, map[string]map[string]string{
		"prod": {"DB_HOST": "localhost", "DB_PASS": "secret"},
	})
	result, err := m.Rotate("prod", map[string]string{"DB_HOST": "DATABASE_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Rotated) != 1 || result.Rotated[0].OldKey != "DB_HOST" || result.Rotated[0].NewKey != "DATABASE_HOST" {
		t.Fatalf("unexpected rotated: %+v", result.Rotated)
	}
	got := store.profiles["prod"]
	if _, exists := got["DB_HOST"]; exists {
		t.Error("old key DB_HOST should have been removed")
	}
	if got["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", got["DATABASE_HOST"])
	}
}

func TestRotateSkipsMissingKeys(t *testing.T) {
	m, _ := newManager(t, map[string]map[string]string{
		"dev": {"API_KEY": "abc123"},
	})
	result, err := m.Rotate("dev", map[string]string{"MISSING_KEY": "NEW_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING_KEY" {
		t.Fatalf("expected MISSING_KEY in skipped, got %+v", result.Skipped)
	}
	if len(result.Rotated) != 0 {
		t.Fatalf("expected no rotations, got %+v", result.Rotated)
	}
}

func TestRotateKeepOldKeyWhenRemoveOldFalse(t *testing.T) {
	store := &fakeStore{profiles: map[string]map[string]string{
		"staging": {"OLD": "value"},
	}}
	m := envrotate.NewManager(store, envrotate.Options{RemoveOld: false})
	_, err := m.Rotate("staging", map[string]string{"OLD": "NEW"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.profiles["staging"]["OLD"] != "value" {
		t.Error("old key should be preserved when RemoveOld=false")
	}
	if store.profiles["staging"]["NEW"] != "value" {
		t.Error("new key should exist")
	}
}

func TestRotateEmptyProfileError(t *testing.T) {
	m, _ := newManager(t, map[string]map[string]string{})
	_, err := m.Rotate("", map[string]string{"A": "B"})
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestRotateEmptyMapError(t *testing.T) {
	m, _ := newManager(t, map[string]map[string]string{"x": {"A": "1"}})
	_, err := m.Rotate("x", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty rotation map")
	}
}

func TestRotateProfileNotFound(t *testing.T) {
	m, _ := newManager(t, map[string]map[string]string{})
	_, err := m.Rotate("ghost", map[string]string{"A": "B"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
