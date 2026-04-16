package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/pin"
	"envoy-cli/internal/store"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "pin-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newManager(t *testing.T) (*pin.Manager, *store.Store) {
	t.Helper()
	s, err := store.Load(filepath.Join(tempDir(t), "store.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Add("production", map[string]string{"ENV": "prod"}); err != nil {
		t.Fatal(err)
	}
	return pin.NewManager(s), s
}

func TestPinAndIsPinned(t *testing.T) {
	m, _ := newManager(t)
	if m.IsPinned("production") {
		t.Fatal("expected not pinned")
	}
	if err := m.Pin("production"); err != nil {
		t.Fatal(err)
	}
	if !m.IsPinned("production") {
		t.Fatal("expected pinned")
	}
}

func TestUnpin(t *testing.T) {
	m, _ := newManager(t)
	m.Pin("production")
	if err := m.Unpin("production"); err != nil {
		t.Fatal(err)
	}
	if m.IsPinned("production") {
		t.Fatal("expected not pinned after unpin")
	}
}

func TestUnpinNotPinned(t *testing.T) {
	m, _ := newManager(t)
	if err := m.Unpin("production"); err == nil {
		t.Fatal("expected error unpinning non-pinned profile")
	}
}

func TestPinNotFound(t *testing.T) {
	m, _ := newManager(t)
	if err := m.Pin("ghost"); err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestListPinned(t *testing.T) {
	m, s := newManager(t)
	s.Add("staging", map[string]string{"ENV": "stage"})
	m.Pin("production")
	pinned, err := m.ListPinned()
	if err != nil {
		t.Fatal(err)
	}
	if len(pinned) != 1 || pinned[0] != "production" {
		t.Fatalf("unexpected pinned list: %v", pinned)
	}
}
