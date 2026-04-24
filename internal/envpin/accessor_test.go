package envpin

import (
	"errors"
	"path/filepath"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{profiles: make(map[string]map[string]string)}
}

func (f *fakeStore) Get(profile string) (map[string]string, error) {
	v, ok := f.profiles[profile]
	if !ok {
		return nil, errors.New("profile not found: " + profile)
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (f *fakeStore) Set(profile string, vars map[string]string) error {
	f.profiles[profile] = vars
	return nil
}

func TestGuardWriteAllowed(t *testing.T) {
	m, _ := NewManager(filepath.Join(tempDir(t), "p.json"))
	store := newFakeStore()
	store.profiles["prod"] = map[string]string{"A": "1", "B": "2"}

	if err := GuardWrite(m, store, "prod", map[string]string{"A": "new", "B": "new"}); err != nil {
		t.Fatalf("GuardWrite: %v", err)
	}
	if store.profiles["prod"]["A"] != "new" {
		t.Error("expected A to be updated")
	}
}

func TestGuardWriteBlocked(t *testing.T) {
	m, _ := NewManager(filepath.Join(tempDir(t), "p.json"))
	_ = m.Pin("prod", "LOCKED")
	store := newFakeStore()
	store.profiles["prod"] = map[string]string{"LOCKED": "original", "FREE": "v"}

	if err := GuardWrite(m, store, "prod", map[string]string{"LOCKED": "changed", "FREE": "new"}); err != nil {
		t.Fatalf("GuardWrite: %v", err)
	}
	if store.profiles["prod"]["LOCKED"] != "original" {
		t.Error("LOCKED should not have been overwritten")
	}
	if store.profiles["prod"]["FREE"] != "new" {
		t.Error("FREE should have been updated")
	}
}

func TestGuardDeleteAllowed(t *testing.T) {
	m, _ := NewManager(filepath.Join(tempDir(t), "p.json"))
	store := newFakeStore()
	store.profiles["dev"] = map[string]string{"A": "1", "B": "2"}

	if err := GuardDelete(m, store, "dev", "A"); err != nil {
		t.Fatalf("GuardDelete: %v", err)
	}
	if _, ok := store.profiles["dev"]["A"]; ok {
		t.Error("expected A to be deleted")
	}
}

func TestGuardDeleteBlocked(t *testing.T) {
	m, _ := NewManager(filepath.Join(tempDir(t), "p.json"))
	_ = m.Pin("dev", "SAFE")
	store := newFakeStore()
	store.profiles["dev"] = map[string]string{"SAFE": "keep"}

	if err := GuardDelete(m, store, "dev", "SAFE"); err == nil {
		t.Error("expected error when deleting pinned key")
	}
}
