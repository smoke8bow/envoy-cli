package quota

import (
	"errors"
	"testing"
)

func newManager() *Manager {
	return NewManager(DefaultPolicy())
}

func TestDefaultMax(t *testing.T)  := newManager()
	if got := m.MaxFor("any"); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}

func TestOverrideMax(t *testing.T) {
	m := newManager()
	m.SetOverride("prod", 10)
	if got := m.MaxFor("prod"); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
	if got := m.MaxFor("dev"); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}

func TestCheckWithinLimit(t *testing.T) {
	m := newManager()
	vars := map[string]string{"A": "1", "B": "2"}
	if err := m.Check("dev", vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckExceedsLimit(t *testing.T) {
	m := newManager()
	m", 2)
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	err := m.Check("small", vars)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestRemoveOverride(t *testing.T) {
	m := newManager()
	m.SetOverride("prod", 5)
	m.RemoveOverride("prod")
	if got := m.MaxFor("prod"); got != 50 {
		t.Fatalf("expected default 50 after removal, got %d", got)
	}
}

func TestCheckExactLimit(t *testing.T) {
	m := newManager()
	m.SetOverride("exact", 3)
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	if err := m.Check("exact", vars); err != nil {
		t.Fatalf("unexpected error at exact limit: %v", err)
	}
}
