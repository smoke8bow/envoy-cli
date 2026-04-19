package alias

import (
	"testing"
)

type fakeStore struct {
	profiles []string
}

func (f *fakeStore) List() ([]string, error) {
	return f.profiles, nil
}

func newManager(profiles ...string) *Manager {
	return NewManager(&fakeStore{profiles: profiles})
}

func TestSetAndResolve(t *testing.T) {
	m := newManager("production")
	if err := m.Set("prod", "production"); err != nil {
		t.Fatal(err)
	}
	got, err := m.Resolve("prod")
	if err != nil {
		t.Fatal(err)
	}
	if got != "production" {
		t.Fatalf("expected production, got %s", got)
	}
}

func TestSetUnknownProfile(t *testing.T) {
	m := newManager()
	if err := m.Set("prod", "production"); err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestSetEmptyAlias(t *testing.T) {
	m := newManager("production")
	if err := m.Set("", "production"); err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestSetInvalidAlias(t *testing.T) {
	m := newManager("production")
	if err := m.Set("my alias!", "production"); err == nil {
		t.Fatal("expected error for invalid alias")
	}
}

func TestResolveNotFound(t *testing.T) {
	m := newManager("production")
	_, err := m.Resolve("unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRemove(t *testing.T) {
	m := newManager("production")
	_ = m.Set("prod", "production")
	if err := m.Remove("prod"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Resolve("prod"); err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveNotFound(t *testing.T) {
	m := newManager()
	if err := m.Remove("ghost"); err == nil {
		t.Fatal("expected error")
	}
}

func TestList(t *testing.T) {
	m := newManager("staging", "production")
	_ = m.Set("prod", "production")
	_ = m.Set("stg", "staging")
	list := m.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(list))
	}
}
