package index

import (
	"errors"
	"testing"
)

// fakeLoader implements ProfileLoader for testing.
type fakeLoader struct {
	profiles map[string]map[string]string
	listErr  error
	getErr   error
}

func (f *fakeLoader) List() ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	names := make([]string, 0, len(f.profiles))
	for n := range f.profiles {
		names = append(names, n)
	}
	return names, nil
}

func (f *fakeLoader) Get(name string) (map[string]string, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.profiles[name], nil
}

func newBuilder(profiles map[string]map[string]string) *Builder {
	return NewBuilder(&fakeLoader{profiles: profiles})
}

func TestBuildEmpty(t *testing.T) {
	b := newBuilder(map[string]map[string]string{})
	idx, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx) != 0 {
		t.Fatalf("expected empty index, got %d entries", len(idx))
	}
}

func TestBuildSingleProfile(t *testing.T) {
	b := newBuilder(map[string]map[string]string{
		"dev": {"HOST": "localhost", "PORT": "8080"},
	})
	idx, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"HOST", "PORT"} {
		profiles := idx.Lookup(key)
		if len(profiles) != 1 || profiles[0] != "dev" {
			t.Errorf("Lookup(%q) = %v, want [dev]", key, profiles)
		}
	}
}

func TestBuildMultipleProfiles(t *testing.T) {
	b := newBuilder(map[string]map[string]string{
		"dev":  {"HOST": "localhost", "DEBUG": "true"},
		"prod": {"HOST": "example.com", "LOG_LEVEL": "warn"},
	})
	idx, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profiles := idx.Lookup("HOST")
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles for HOST, got %d", len(profiles))
	}
	if profiles[0] != "dev" || profiles[1] != "prod" {
		t.Errorf("unexpected order: %v", profiles)
	}
}

func TestLookupMissingKey(t *testing.T) {
	idx := Index{}
	if got := idx.Lookup("MISSING"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestEntriesSortedByKey(t *testing.T) {
	b := newBuilder(map[string]map[string]string{
		"dev": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	})
	idx, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := idx.Entries()
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("entry[%d].Key = %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestBuildListError(t *testing.T) {
	loader := &fakeLoader{listErr: errors.New("storage failure")}
	_, err := NewBuilder(loader).Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
