package readonly

import (
	"testing"
)

func newManager() *Manager { return NewManager() }

func TestSetAndIsReadOnly(t *testing.T) {
	m := newManager()
	if err := m.Set("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.IsReadOnly("prod") {
		t.Fatal("expected prod to be read-only")
	}
}

func TestSetEmptyName(t *testing.T) {
	m := newManager()
	if err := m.Set(""); err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestUnset(t *testing.T) {
	m := newManager()
	_ = m.Set("prod")
	if err := m.Unset("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.IsReadOnly("prod") {
		t.Fatal("expected prod to no longer be read-only")
	}
}

func TestUnsetNotReadOnly(t *testing.T) {
	m := newManager()
	if err := m.Unset("staging"); err == nil {
		t.Fatal("expected error when unsetting a non-read-only profile")
	}
}

func TestCheckAllowsMutable(t *testing.T) {
	m := newManager()
	if err := m.Check("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckBlocksReadOnly(t *testing.T) {
	m := newManager()
	_ = m.Set("prod")
	if err := m.Check("prod"); err == nil {
		t.Fatal("expected error for read-only profile")
	}
}

func TestList(t *testing.T) {
	m := newManager()
	_ = m.Set("prod")
	_ = m.Set("staging")
	list := m.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
}
