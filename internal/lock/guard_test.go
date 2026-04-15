package lock_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-cli/internal/lock"
)

func newGuard(t *testing.T) (*lock.Manager, *lock.Guard) {
	t.Helper()
	m := lock.NewManager(tempDir(t))
	return m, lock.NewGuard(m)
}

func TestRequireUnlocked(t *testing.T) {
	_, g := newGuard(t)
	if err := g.Require("prod"); err != nil {
		t.Fatalf("expected no error for unlocked profile, got %v", err)
	}
}

func TestRequireLocked(t *testing.T) {
	m, g := newGuard(t)
	_ = m.Lock("prod")
	if err := g.Require("prod"); err == nil {
		t.Fatal("expected error for locked profile")
	}
}

func TestWithUnlockedRuns(t *testing.T) {
	_, g := newGuard(t)
	called := false
	err := g.WithUnlocked("dev", func() error {
		called = true
		return nil
	})
	if err != nil || !called {
		t.Fatalf("expected fn to run, err=%v called=%v", err, called)
	}
}

func TestWithUnlockedBlocked(t *testing.T) {
	m, g := newGuard(t)
	_ = m.Lock("dev")
	called := false
	err := g.WithUnlocked("dev", func() error {
		called = true
		return nil
	})
	if err == nil {
		t.Fatal("expected error when profile is locked")
	}
	if called {
		t.Fatal("fn should not have been called")
	}
}

func TestStatusUnlocked(t *testing.T) {
	_, g := newGuard(t)
	s := g.Status("staging")
	if !strings.Contains(s, "unlocked") {
		t.Fatalf("unexpected status: %s", s)
	}
}

func TestStatusLocked(t *testing.T) {
	m, g := newGuard(t)
	_ = m.Lock("staging")
	s := g.Status("staging")
	if !strings.Contains(s, "locked") {
		t.Fatalf("unexpected status: %s", s)
	}
}
