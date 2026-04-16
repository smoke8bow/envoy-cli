package inherit_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envoy-cli/internal/inherit"
)

type fakeStore struct {
	data map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.data[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	f.data[name] = vars
	return nil
}

func newInheritor() (*inherit.Inheritor, *fakeStore) {
	s := &fakeStore{data: map[string]map[string]string{
		"base": {"HOST": "localhost", "PORT": "5432", "DEBUG": "false"},
		"prod": {"HOST": "prod.example.com", "LOG_LEVEL": "warn"},
	}}
	return inherit.NewInheritor(s), s
}

func TestApplyChildKeepsOwnValues(t *testing.T) {
	inh, _ := newInheritor()
	result, err := inh.Apply("base", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if result["HOST"] != "prod.example.com" {
		t.Errorf("expected child HOST to win, got %s", result["HOST"])
	}
}

func TestApplyMissingKeysFilledFromParent(t *testing.T) {
	inh, _ := newInheritor()
	result, err := inh.Apply("base", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT from parent, got %s", result["PORT"])
	}
	if result["DEBUG"] != "false" {
		t.Errorf("expected DEBUG from parent, got %s", result["DEBUG"])
	}
}

func TestApplyChildOnlyKeyPreserved(t *testing.T) {
	inh, _ := newInheritor()
	result, err := inh.Apply("base", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if result["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL from child, got %s", result["LOG_LEVEL"])
	}
}

func TestApplyParentNotFound(t *testing.T) {
	inh, _ := newInheritor()
	_, err := inh.Apply("missing", "prod")
	if err == nil {
		t.Fatal("expected error for missing parent")
	}
}

func TestCommitPersists(t *testing.T) {
	inh, s := newInheritor()
	_, err := inh.Commit("base", "prod")
	if err != nil {
		t.Fatal(err)
	}
	if s.data["prod"]["PORT"] != "5432" {
		t.Errorf("expected PORT persisted, got %s", s.data["prod"]["PORT"])
	}
}
