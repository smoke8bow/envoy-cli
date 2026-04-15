package lock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/lock"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "lock-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return filepath.Join(dir, "locks")
}

func TestLockAndIsLocked(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	if m.IsLocked("prod") {
		t.Fatal("expected profile to not be locked initially")
	}
	if err := m.Lock("prod"); err != nil {
		t.Fatalf("Lock() error: %v", err)
	}
	if !m.IsLocked("prod") {
		t.Fatal("expected profile to be locked after Lock()")
	}
}

func TestLockDuplicate(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	if err := m.Lock("staging"); err != nil {
		t.Fatalf("first Lock() error: %v", err)
	}
	if err := m.Lock("staging"); err != lock.ErrLocked {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestUnlock(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	_ = m.Lock("dev")
	if err := m.Unlock("dev"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}
	if m.IsLocked("dev") {
		t.Fatal("expected profile to be unlocked after Unlock()")
	}
}

func TestUnlockNotLocked(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	if err := m.Unlock("ghost"); err != lock.ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestLockedAt(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	if !m.LockedAt("x").IsZero() {
		t.Fatal("expected zero time for unlocked profile")
	}
	_ = m.Lock("x")
	if m.LockedAt("x").IsZero() {
		t.Fatal("expected non-zero time after locking")
	}
}

func TestList(t *testing.T) {
	m := lock.NewManager(tempDir(t))
	names, err := m.List()
	if err != nil || len(names) != 0 {
		t.Fatalf("expected empty list, got %v %v", names, err)
	}
	_ = m.Lock("alpha")
	_ = m.Lock("beta")
	names, err = m.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 locked profiles, got %d", len(names))
	}
}
