package chain_test

import (
	"errors"
	"testing"

	"github.com/envoy-cli/envoy/internal/chain"
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

func TestFromProfilesSuccess(t *testing.T) {
	store := &fakeStore{
		profiles: map[string]map[string]string{
			"base": {"X": "1"},
			"prod": {"X": "2", "Y": "3"},
		},
	}
	c, err := chain.FromProfiles(store, []string{"base", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := c.Compose()
	if r.Vars["X"] != "2" {
		t.Fatalf("expected X=2, got %q", r.Vars["X"])
	}
	if r.Source["X"] != "prod" {
		t.Fatalf("expected source prod, got %q", r.Source["X"])
	}
}

func TestFromProfilesNotFound(t *testing.T) {
	store := &fakeStore{profiles: map[string]map[string]string{}}
	_, err := chain.FromProfiles(store, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestFromProfilesEmpty(t *testing.T) {
	store := &fakeStore{profiles: map[string]map[string]string{}}
	_, err := chain.FromProfiles(store, nil)
	if err == nil {
		t.Fatal("expected error for empty names")
	}
}
