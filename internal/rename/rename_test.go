package rename_test

import (
	"errors"
	"testing"

	"envoy-cli/internal/rename"
)

// fakeStore is an in-memory Store implementation for testing.
type fakeStore struct {
	profiles map[string]map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{profiles: make(map[string]map[string]string)}
}

func (f *fakeStore) List() []string {
	names := make([]string, 0, len(f.profiles))
	for k := range f.profiles {
		names = append(names, k)
	}
	return names
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
	f.profiles[name] = vars
	return nil
}

func (f *fakeStore) Delete(name string) error {
	delete(f.profiles, name)
	return nil
}

func newRenamer() (*rename.Renamer, *fakeStore) {
	s := newFakeStore()
	return rename.NewRenamer(s), s
}

func TestRenameSuccess(t *testing.T) {
	r, s := newRenamer()
	_ = s.Set("dev", map[string]string{"FOO": "bar"})

	if err := r.Rename("dev", "development"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := s.Get("development"); err != nil {
		t.Error("expected 'development' profile to exist")
	}
	if _, err := s.Get("dev"); err == nil {
		t.Error("expected 'dev' profile to be removed")
	}
}

func TestRenamePreservesVars(t *testing.T) {
	r, s := newRenamer()
	_ = s.Set("staging", map[string]string{"API_URL": "https://staging.example.com", "DEBUG": "true"})

	if err := r.Rename("staging", "stg"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, err := s.Get("stg")
	if err != nil {
		t.Fatal("expected 'stg' profile to exist")
	}
	if vars["API_URL"] != "https://staging.example.com" {
		t.Errorf("expected API_URL to be preserved, got %q", vars["API_URL"])
	}
}

func TestRenameNotFound(t *testing.T) {
	r, _ := newRenamer()
	err := r.Rename("ghost", "newname")
	if !errors.Is(err, rename.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRenameAlreadyExists(t *testing.T) {
	r, s := newRenamer()
	_ = s.Set("a", map[string]string{"X": "1"})
	_ = s.Set("b", map[string]string{"Y": "2"})

	err := r.Rename("a", "b")
	if !errors.Is(err, rename.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestRenameSameName(t *testing.T) {
	r, s := newRenamer()
	_ = s.Set("prod", map[string]string{})

	err := r.Rename("prod", "prod")
	if !errors.Is(err, rename.ErrSameName) {
		t.Errorf("expected ErrSameName, got %v", err)
	}
}
