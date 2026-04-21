package draft_test

import (
	"errors"
	"testing"

	"github.com/user/envoy-cli/internal/draft"
)

func newManager() *draft.Manager { return draft.NewManager() }

func TestCreateAndGet(t *testing.T) {
	m := newManager()
	if err := m.Create("staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vars, err := m.Get("staging")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(vars) != 0 {
		t.Errorf("expected empty vars, got %v", vars)
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := newManager()
	_ = m.Create("dup")
	if err := m.Create("dup"); err == nil {
		t.Error("expected error for duplicate draft")
	}
}

func TestSetAndGet(t *testing.T) {
	m := newManager()
	_ = m.Create("d")
	_ = m.Set("d", "FOO", "bar")
	_ = m.Set("d", "BAZ", "qux")
	vars, _ := m.Get("d")
	if vars["FOO"] != "bar" || vars["BAZ"] != "qux" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestSetNotFound(t *testing.T) {
	m := newManager()
	err := m.Set("missing", "K", "V")
	if !errors.Is(err, draft.ErrNoDraft) {
		t.Errorf("expected ErrNoDraft, got %v", err)
	}
}

func TestDeleteKey(t *testing.T) {
	m := newManager()
	_ = m.Create("d")
	_ = m.Set("d", "FOO", "bar")
	_ = m.Delete("d", "FOO")
	vars, _ := m.Get("d")
	if _, ok := vars["FOO"]; ok {
		t.Error("expected FOO to be deleted")
	}
}

func TestDiscard(t *testing.T) {
	m := newManager()
	_ = m.Create("d")
	_ = m.Discard("d")
	_, err := m.Get("d")
	if !errors.Is(err, draft.ErrNoDraft) {
		t.Errorf("expected ErrNoDraft after discard, got %v", err)
	}
}

func TestDiscardNotFound(t *testing.T) {
	m := newManager()
	if err := m.Discard("ghost"); !errors.Is(err, draft.ErrNoDraft) {
		t.Errorf("expected ErrNoDraft, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := newManager()
	_ = m.Create("a")
	_ = m.Create("b")
	names := m.List()
	if len(names) != 2 {
		t.Errorf("expected 2 drafts, got %d", len(names))
	}
}
