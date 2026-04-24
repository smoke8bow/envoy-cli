package envpin

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "envpin-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func newManager(t *testing.T) *Manager {
	t.Helper()
	m, err := NewManager(filepath.Join(tempDir(t), "envpin.json"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestPinAndIsPinned(t *testing.T) {
	m := newManager(t)
	if err := m.Pin("prod", "DB_PASS"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	if !m.IsPinned("prod", "DB_PASS") {
		t.Error("expected DB_PASS to be pinned")
	}
	if m.IsPinned("prod", "OTHER") {
		t.Error("OTHER should not be pinned")
	}
}

func TestPinDuplicate(t *testing.T) {
	m := newManager(t)
	if err := m.Pin("prod", "KEY"); err != nil {
		t.Fatalf("first Pin: %v", err)
	}
	if err := m.Pin("prod", "KEY"); err == nil {
		t.Error("expected error on duplicate pin")
	}
}

func TestUnpin(t *testing.T) {
	m := newManager(t)
	_ = m.Pin("dev", "SECRET")
	if err := m.Unpin("dev", "SECRET"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}
	if m.IsPinned("dev", "SECRET") {
		t.Error("expected SECRET to be unpinned")
	}
}

func TestUnpinNotPinned(t *testing.T) {
	m := newManager(t)
	if err := m.Unpin("dev", "MISSING"); err == nil {
		t.Error("expected error when unpinning non-pinned key")
	}
}

func TestKeys(t *testing.T) {
	m := newManager(t)
	_ = m.Pin("prod", "Z_KEY")
	_ = m.Pin("prod", "A_KEY")
	keys := m.Keys("prod")
	if len(keys) != 2 || keys[0] != "A_KEY" || keys[1] != "Z_KEY" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestFilterWritable(t *testing.T) {
	m := newManager(t)
	_ = m.Pin("prod", "LOCKED")
	vars := map[string]string{"LOCKED": "x", "FREE": "y"}
	out := m.FilterWritable("prod", vars)
	if _, ok := out["LOCKED"]; ok {
		t.Error("LOCKED should have been filtered out")
	}
	if out["FREE"] != "y" {
		t.Error("FREE should be present")
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "envpin.json")
	m1, _ := NewManager(path)
	_ = m1.Pin("staging", "API_KEY")

	m2, err := NewManager(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !m2.IsPinned("staging", "API_KEY") {
		t.Error("expected API_KEY to persist across reload")
	}
}
