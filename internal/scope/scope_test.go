package scope

import (
	"testing"
)

func newManager() *Manager {
	return NewManager()
}

func TestCreateAndGet(t *testing.T) {
	m := newManager()
	s, err := m.Create("prod", map[string]string{"env": "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "prod" {
		t.Errorf("expected name prod, got %s", s.Name)
	}
	got, err := m.Get("prod")
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	if got.Labels["env"] != "production" {
		t.Errorf("expected label env=production")
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := newManager()
	m.Create("dev", nil)
	_, err := m.Create("dev", nil)
	if err == nil {
		t.Error("expected error for duplicate scope")
	}
}

func TestCreateEmptyName(t *testing.T) {
	m := newManager()
	_, err := m.Create("", nil)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestGetNotFound(t *testing.T) {
	m := newManager()
	_, err := m.Get("missing")
	if err == nil {
		t.Error("expected error for missing scope")
	}
}

func TestDelete(t *testing.T) {
	m := newManager()
	m.Create("staging", nil)
	if err := m.Delete("staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := m.Get("staging"); err == nil {
		t.Error("expected scope to be deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newManager()
	if err := m.Delete("ghost"); err == nil {
		t.Error("expected error deleting non-existent scope")
	}
}

func TestList(t *testing.T) {
	m := newManager()
	m.Create("a", nil)
	m.Create("b", nil)
	if len(m.List()) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(m.List()))
	}
}

func TestMatch(t *testing.T) {
	s := &Scope{Name: "x", Labels: map[string]string{"team": "backend", "env": "prod"}}
	if !s.Match(map[string]string{"team": "backend"}) {
		t.Error("expected match")
	}
	if s.Match(map[string]string{"team": "frontend"}) {
		t.Error("expected no match")
	}
}
