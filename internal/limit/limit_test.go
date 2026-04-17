package limit

import (
	"errors"
	"testing"
)

func newManager() *Manager { return NewManager(10) }

func TestDefaultLimit(t *testing.T) {
	m := newManager()
	if got := m.GetLimit("any"); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestOverrideLimit(t *testing.T) {
	m := newManager()
	m.SetLimit("prod", 3)
	if got := m.GetLimit("prod"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
	if got := m.GetLimit("dev"); got != 10 {
		t.Fatalf("expected default 10, got %d", got)
	}
}

func TestCheckWithinLimit(t *testing.T) {
	m := newManager()
	if err := m.Check("dev", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckExceedsLimit(t *testing.T) {
	m := newManager()
	err := m.Check("dev", 11)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestCheckVars(t *testing.T) {
	m := newManager()
	m.SetLimit("small", 2)
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	err := m.CheckVars("small", vars)
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestCheckVarsWithinLimit(t *testing.T) {
	m := newManager()
	vars := map[string]string{"X": "1"}
	if err := m.CheckVars("any", vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
