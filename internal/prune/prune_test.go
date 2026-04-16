package prune

import (
	"errors"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
	failDelete string
}

func (f *fakeStore) List() ([]string, error) {
	names := make([]string, 0, len(f.profiles))
	for k := range f.profiles {
		names = append(names, k)
	}
	return names, nil
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (f *fakeStore) Delete(name string) error {
	if f.failDelete == name {
		return errors.New("delete failed")
	}
	delete(f.profiles, name)
	return nil
}

func newManager(t *testing.T) (*Manager, *fakeStore) {
	t.Helper()
	fs := &fakeStore{profiles: map[string]map[string]string{}}
	return NewManager(fs), fs
}

func TestDryRunNoCandidates(t *testing.T) {
	m, fs := newManager(t)
	fs.profiles["prod"] = map[string]string{"KEY": "val"}
	candidates, err := m.DryRun()
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 0 {
		t.Fatalf("expected 0 candidates, got %d", len(candidates))
	}
}

func TestDryRunFindsEmpty(t *testing.T) {
	m, fs := newManager(t)
	fs.profiles["empty"] = map[string]string{}
	fs.profiles["full"] = map[string]string{"A": "1"}
	candidates, err := m.DryRun()
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0] != "empty" {
		t.Fatalf("expected [empty], got %v", candidates)
	}
}

func TestRunRemovesEmpty(t *testing.T) {
	m, fs := newManager(t)
	fs.profiles["ghost"] = map[string]string{}
	fs.profiles["keep"] = map[string]string{"X": "y"}
	res, err := m.Run()
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "ghost" {
		t.Fatalf("unexpected removed: %v", res.Removed)
	}
	if _, ok := fs.profiles["ghost"]; ok {
		t.Fatal("ghost profile should have been deleted")
	}
}

func TestRunSkipsOnDeleteError(t *testing.T) {
	m, fs := newManager(t)
	fs.profiles["bad"] = map[string]string{}
	fs.failDelete = "bad"
	res, err := m.Run()
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "bad" {
		t.Fatalf("expected bad in skipped, got %v", res.Skipped)
	}
}
