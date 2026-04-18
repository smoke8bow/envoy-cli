package prefix

import (
	"errors"
	"testing"
)

type fakeStore struct {
	data map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.data[name]
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
	f.data[name] = vars
	return nil
}

func TestApplyToProfile(t *testing.T) {
	fs := &fakeStore{data: map[string]map[string]string{
		"dev": {"HOST": "localhost", "PORT": "5432"},
	}}
	if err := ApplyToProfile(fs, "dev", "PG_"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.data["dev"]["PG_HOST"] != "localhost" {
		t.Fatalf("expected PG_HOST, got %v", fs.data["dev"])
	}
}

func TestApplyToProfileNotFound(t *testing.T) {
	fs := &fakeStore{data: map[string]map[string]string{}}
	if err := ApplyToProfile(fs, "missing", "X_"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestStripFromProfile(t *testing.T) {
	fs := &fakeStore{data: map[string]map[string]string{
		"dev": {"PG_HOST": "localhost", "PG_PORT": "5432"},
	}}
	if err := StripFromProfile(fs, "dev", "PG_", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.data["dev"]["HOST"] != "localhost" {
		t.Fatalf("expected HOST key after strip, got %v", fs.data["dev"])
	}
}
