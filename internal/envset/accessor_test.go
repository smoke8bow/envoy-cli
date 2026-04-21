package envset_test

import (
	"errors"
	"testing"

	"envoy-cli/internal/envset"
)

type fakeStore struct {
	profiles map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	p, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return p, nil
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		profiles: map[string]map[string]string{
			"base": {"A": "1", "B": "2"},
			"override": {"B": "99", "C": "3"},
		},
	}
}

func TestApplyToProfilesUnion(t *testing.T) {
	store := newFakeStore()
	got, err := envset.ApplyToProfiles(store, envset.OpUnion, "base", "override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 keys, got %d", len(got))
	}
	if got["B"] != "99" {
		t.Errorf("expected B=99 from override, got %s", got["B"])
	}
}

func TestApplyToProfilesNotFound(t *testing.T) {
	store := newFakeStore()
	_, err := envset.ApplyToProfiles(store, envset.OpUnion, "missing", "override")
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestApplyToProfilesInvalidOp(t *testing.T) {
	store := newFakeStore()
	_, err := envset.ApplyToProfiles(store, "bad", "base", "override")
	if err == nil {
		t.Error("expected error for invalid op")
	}
}

func TestApplyToProfilesDifference(t *testing.T) {
	store := newFakeStore()
	got, err := envset.ApplyToProfiles(store, envset.OpDifference, "base", "override")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// base has A,B; override has B,C → difference = A only
	if len(got) != 1 || got["A"] != "1" {
		t.Errorf("unexpected difference result: %v", got)
	}
}
