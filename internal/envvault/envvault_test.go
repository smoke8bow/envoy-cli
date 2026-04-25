package envvault_test

import (
	"errors"
	"testing"

	"github.com/user/envoy-cli/internal/envvault"
)

const testPass = "supersecret"

// fakeStorage is an in-memory Storage for tests.
type fakeStorage struct {
	blobs map[string]string
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{blobs: map[string]string{}}
}

func (f *fakeStorage) Load() (map[string]string, error) {
	copy := make(map[string]string, len(f.blobs))
	for k, v := range f.blobs {
		copy[k] = v
	}
	return copy, nil
}

func (f *fakeStorage) Save(blobs map[string]string) error {
	f.blobs = blobs
	return nil
}

func newManager(t *testing.T) *envvault.Manager {
	t.Helper()
	return envvault.NewManager(newFakeStorage(), testPass)
}

func TestPutAndGet(t *testing.T) {
	m := newManager(t)
	vars := map[string]string{"KEY": "value", "FOO": "bar"}
	if err := m.Put("prod", vars); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := m.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got["KEY"] != "value" || got["FOO"] != "bar" {
		t.Errorf("unexpected vars: %v", got)
	}
}

func TestPutDuplicate(t *testing.T) {
	m := newManager(t)
	_ = m.Put("prod", map[string]string{"A": "1"})
	err := m.Put("prod", map[string]string{"B": "2"})
	if !errors.Is(err, envvault.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestSetOverwrites(t *testing.T) {
	m := newManager(t)
	_ = m.Put("prod", map[string]string{"A": "old"})
	_ = m.Set("prod", map[string]string{"A": "new"})
	got, _ := m.Get("prod")
	if got["A"] != "new" {
		t.Errorf("expected 'new', got %q", got["A"])
	}
}

func TestGetNotFound(t *testing.T) {
	m := newManager(t)
	_, err := m.Get("missing")
	if !errors.Is(err, envvault.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := newManager(t)
	_ = m.Put("prod", map[string]string{"X": "1"})
	if err := m.Delete("prod"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("prod")
	if !errors.Is(err, envvault.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newManager(t)
	err := m.Delete("ghost")
	if !errors.Is(err, envvault.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := newManager(t)
	_ = m.Put("alpha", map[string]string{"A": "1"})
	_ = m.Put("beta", map[string]string{"B": "2"})
	names, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
