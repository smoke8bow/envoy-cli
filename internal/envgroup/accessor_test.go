package envgroup

import (
	"errors"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func newFakeStore() *fakeStore {
	return &fakeStore{profiles: map[string]map[string]string{}}
}

func TestGroupProfileSuccess(t *testing.T) {
	store := newFakeStore()
	store.profiles["prod"] = map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
		"APP_PORT": "8080",
	}
	r, err := GroupProfile(store, "prod", DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(r.Groups))
	}
}

func TestGroupProfileNotFound(t *testing.T) {
	store := newFakeStore()
	_, err := GroupProfile(store, "missing", DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestGroupProfileEmpty(t *testing.T) {
	store := newFakeStore()
	store.profiles["empty"] = map[string]string{}
	r, err := GroupProfile(store, "empty", DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups) != 0 {
		t.Errorf("expected no groups for empty profile")
	}
}
