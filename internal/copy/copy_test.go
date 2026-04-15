package copy

import (
	"errors"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{profiles: make(map[string]map[string]string)}
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

func (f *fakeStore) Create(name string, vars map[string]string) error {
	f.profiles[name] = vars
	return nil
}

func (f *fakeStore) Exists(name string) bool {
	_, ok := f.profiles[name]
	return ok
}

func newCopier() (*Copier, *fakeStore) {
	s := newFakeStore()
	return NewCopier(s), s
}

func TestCopyKeysSuccess(t *testing.T) {
	c, s := newCopier()
	s.profiles["src"] = map[string]string{"A": "1", "B": "2", "C": "3"}

	if err := c.CopyKeys("src", "dst", []string{"A", "C"}, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dst := s.profiles["dst"]
	if dst["A"] != "1" || dst["C"] != "3" {
		t.Errorf("expected A=1 C=3, got %v", dst)
	}
	if _, ok := dst["B"]; ok {
		t.Error("key B should not be present in dst")
	}
}

func TestCopyKeysNoOverwrite(t *testing.T) {
	c, s := newCopier()
	s.profiles["src"] = map[string]string{"A": "new"}
	s.profiles["dst"] = map[string]string{"A": "old"}

	if err := c.CopyKeys("src", "dst", []string{"A"}, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.profiles["dst"]["A"] != "old" {
		t.Errorf("expected A to remain 'old', got %q", s.profiles["dst"]["A"])
	}
}

func TestCopyKeysOverwrite(t *testing.T) {
	c, s := newCopier()
	s.profiles["src"] = map[string]string{"A": "new"}
	s.profiles["dst"] = map[string]string{"A": "old"}

	if err := c.CopyKeys("src", "dst", []string{"A"}, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.profiles["dst"]["A"] != "new" {
		t.Errorf("expected A='new', got %q", s.profiles["dst"]["A"])
	}
}

func TestCopyKeysMissingKey(t *testing.T) {
	c, s := newCopier()
	s.profiles["src"] = map[string]string{"A": "1"}

	err := c.CopyKeys("src", "dst", []string{"MISSING"}, false)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestCopyKeysSourceNotFound(t *testing.T) {
	c, _ := newCopier()
	err := c.CopyKeys("ghost", "dst", []string{"A"}, false)
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}

func TestCopyKeysEmptyNames(t *testing.T) {
	c, _ := newCopier()
	if err := c.CopyKeys("", "dst", []string{"A"}, false); err == nil {
		t.Error("expected error for empty src")
	}
	if err := c.CopyKeys("src", "", []string{"A"}, false); err == nil {
		t.Error("expected error for empty dst")
	}
}

func TestCopyKeysEmptyKeys(t *testing.T) {
	c, s := newCopier()
	s.profiles["src"] = map[string]string{"A": "1"}
	if err := c.CopyKeys("src", "dst", []string{}, false); err == nil {
		t.Error("expected error for empty keys slice")
	}
}
