package supersede_test

import (
	"errors"
	"sort"
	"testing"

	"envoy-cli/internal/supersede"
)

// fakeStore is an in-memory store for testing.
type fakeStore struct {
	data    map[string]map[string]string
	getErr  error
	setErr  error
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	v, ok := f.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.data[name] = vars
	return nil
}

func newManager(profiles map[string]map[string]string) *supersede.Manager {
	return supersede.NewManager(&fakeStore{data: profiles})
}

func TestApplyAllKeys(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"A": "1", "B": "2"},
		"src":  {"A": "10", "B": "20", "C": "30"},
	})
	applied, err := m.Apply("base", "src", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(applied) != 3 {
		t.Errorf("expected 3 applied keys, got %d", len(applied))
	}
}

func TestApplySelectedKeys(t *testing.T) {
	store := &fakeStore{data: map[string]map[string]string{
		"base": {"A": "1", "B": "2"},
		"src":  {"A": "99", "B": "88"},
	}}
	m := supersede.NewManager(store)
	applied, err := m.Apply("base", "src", []string{"A"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(applied) != 1 || applied[0] != "A" {
		t.Errorf("expected [A], got %v", applied)
	}
	result := store.data["base"]
	if result["A"] != "99" {
		t.Errorf("expected A=99, got %s", result["A"])
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2 (unchanged), got %s", result["B"])
	}
}

func TestApplySkipsMissingKeys(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"A": "1"},
		"src":  {"B": "2"},
	})
	applied, err := m.Apply("base", "src", []string{"MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(applied) != 0 {
		t.Errorf("expected no applied keys, got %v", applied)
	}
}

func TestApplyDoesNotMutateSrc(t *testing.T) {
	store := &fakeStore{data: map[string]map[string]string{
		"base": {"X": "orig"},
		"src":  {"X": "new", "Y": "extra"},
	}}
	m := supersede.NewManager(store)
	_, err := m.Apply("base", "src", []string{"X"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.data["src"]["Y"] != "extra" {
		t.Error("source profile was mutated")
	}
}

func TestApplyEmptyDestinationError(t *testing.T) {
	m := newManager(map[string]map[string]string{})
	_, err := m.Apply("", "src", nil)
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestApplyEmptySourceError(t *testing.T) {
	m := newManager(map[string]map[string]string{})
	_, err := m.Apply("base", "", nil)
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestApplyReturnsSortedKeys(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"base": {"A": "1", "B": "2", "C": "3"},
		"src":  {"A": "9", "B": "8", "C": "7"},
	})
	applied, err := m.Apply("base", "src", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sort.StringsAreSorted(applied) {
		t.Errorf("expected sorted applied keys, got %v", applied)
	}
}
