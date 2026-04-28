package envpromote

import (
	"errors"
	"testing"
)

type fakeStore struct {
	data map[string]map[string]string
}

func newFakeStore(profiles map[string]map[string]string) *fakeStore {
	return &fakeStore{data: profiles}
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.data[name]
	if !ok {
		return nil, errors.New("profile not found: " + name)
	}
	return v, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	f.data[name] = vars
	return nil
}

func newManager(t *testing.T) (*Manager, *fakeStore) {
	t.Helper()
	store := newFakeStore(map[string]map[string]string{
		"dev":  {"APP_ENV": "development", "DB_HOST": "localhost", "DEBUG": "true"},
		"prod": {"APP_ENV": "production"},
	})
	return NewManager(store), store
}

func TestPromoteAllKeys(t *testing.T) {
	m, store := newManager(t)
	result, err := m.Promote("dev", "prod", DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
	if result["APP_ENV"] != "development" {
		t.Errorf("expected APP_ENV overwritten to development, got %q", result["APP_ENV"])
	}
	if store.data["prod"]["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true in persisted prod")
	}
}

func TestPromoteSelectedKeys(t *testing.T) {
	m, _ := newManager(t)
	opts := Options{Keys: []string{"DB_HOST"}, Overwrite: true}
	result, err := m.Promote("dev", "prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["DEBUG"]; ok {
		t.Error("DEBUG should not have been promoted")
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
}

func TestPromoteNoOverwrite(t *testing.T) {
	m, _ := newManager(t)
	opts := Options{Overwrite: false}
	result, err := m.Promote("dev", "prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// prod already had APP_ENV=production; should not be overwritten.
	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production (no overwrite), got %q", result["APP_ENV"])
	}
}

func TestPromoteSrcNotFound(t *testing.T) {
	m, _ := newManager(t)
	_, err := m.Promote("missing", "prod", DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing src profile")
	}
}

func TestPromoteSameSrcDst(t *testing.T) {
	m, _ := newManager(t)
	_, err := m.Promote("dev", "dev", DefaultOptions())
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestPromoteEmptySrc(t *testing.T) {
	m, _ := newManager(t)
	_, err := m.Promote("", "prod", DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}
