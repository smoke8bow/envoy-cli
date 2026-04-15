package snapshot

import (
	"errors"
	"strings"
	"testing"
)

// fakeStore is an in-memory snapshotStore for tests.
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

func (f *fakeStore) Save(name string, vars map[string]string) error {
	f.data[name] = vars
	return nil
}

func TestTakeCreatesSnapshot(t *testing.T) {
	fs := newFakeStore()
	fs.data["prod"] = map[string]string{"APP_ENV": "production", "PORT": "8080"}

	m := NewManager(fs)
	entry, err := m.Take("prod", "before deploy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry.Profile != "prod" {
		t.Errorf("expected profile=prod, got %s", entry.Profile)
	}
	if !strings.HasPrefix(entry.ID, "prod__snap_") {
		t.Errorf("unexpected snapshot ID format: %s", entry.ID)
	}
	if entry.Note != "before deploy" {
		t.Errorf("expected note='before deploy', got %s", entry.Note)
	}
	if _, ok := fs.data[entry.ID]; !ok {
		t.Error("snapshot not persisted in store")
	}
}

func TestTakeProfileNotFound(t *testing.T) {
	fs := newFakeStore()
	m := NewManager(fs)
	_, err := m.Take("missing", "")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestTakeDoesNotMutateOriginal(t *testing.T) {
	fs := newFakeStore()
	fs.data["dev"] = map[string]string{"KEY": "original"}

	m := NewManager(fs)
	entry, _ := m.Take("dev", "")

	// Mutate snapshot vars directly
	entry.Vars["KEY"] = "mutated"

	if fs.data["dev"]["KEY"] != "original" {
		t.Error("Take mutated the original profile vars")
	}
}

func TestRestoreSuccess(t *testing.T) {
	fs := newFakeStore()
	fs.data["prod"] = map[string]string{"APP_ENV": "production"}

	m := NewManager(fs)
	entry, _ := m.Take("prod", "")

	// Overwrite prod
	fs.data["prod"] = map[string]string{"APP_ENV": "staging"}

	if err := m.Restore(entry.ID, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.data["prod"]["APP_ENV"] != "production" {
		t.Errorf("restore failed: got %s", fs.data["prod"]["APP_ENV"])
	}
}

func TestRestoreSnapshotNotFound(t *testing.T) {
	fs := newFakeStore()
	m := NewManager(fs)
	err := m.Restore("nonexistent__snap_000", "prod")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
