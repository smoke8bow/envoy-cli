package freeze_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/freeze"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "freeze-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func newManager(t *testing.T) *freeze.Manager {
	t.Helper()
	m, err := freeze.NewManager(filepath.Join(tempDir(t), "frozen.json"))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestFreezeAndIsFrozen(t *testing.T) {
	m := newManager(t)
	if m.IsFrozen("prod") {
		t.Fatal("expected not frozen")
	}
	if err := m.Freeze("prod"); err != nil {
		t.Fatal(err)
	}
	if !m.IsFrozen("prod") {
		t.Fatal("expected frozen")
	}
}

func TestUnfreeze(t *testing.T) {
	m := newManager(t)
	_ = m.Freeze("staging")
	if err := m.Unfreeze("staging"); err != nil {
		t.Fatal(err)
	}
	if m.IsFrozen("staging") {
		t.Fatal("expected not frozen after unfreeze")
	}
}

func TestUnfreezeNotFrozen(t *testing.T) {
	m := newManager(t)
	if err := m.Unfreeze("dev"); err != freeze.ErrNotFrozen {
		t.Fatalf("expected ErrNotFrozen, got %v", err)
	}
}

func TestFreezeEmptyNameError(t *testing.T) {
	m := newManager(t)
	if err := m.Freeze(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestList(t *testing.T) {
	m := newManager(t)
	_ = m.Freeze("a")
	_ = m.Freeze("b")
	list := m.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 frozen profiles, got %d", len(list))
	}
}

func TestGuard(t *testing.T) {
	m := newManager(t)
	if err := m.Guard("x"); err != nil {
		t.Fatal("expected no error for unfrozen profile")
	}
	_ = m.Freeze("x")
	if err := m.Guard("x"); err != freeze.ErrFrozen {
		t.Fatalf("expected ErrFrozen, got %v", err)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "frozen.json")

	m1, _ := freeze.NewManager(path)
	_ = m1.Freeze("prod")

	m2, err := freeze.NewManager(path)
	if err != nil {
		t.Fatal(err)
	}
	if !m2.IsFrozen("prod") {
		t.Fatal("expected frozen state to persist")
	}
}
