package envclip

import (
	"testing"
)

func newManager() *Manager { return NewManager() }

func TestCopyAllKeys(t *testing.T) {
	m := newManager()
	vars := map[string]string{"A": "1", "B": "2"}
	if err := m.Copy("dev", vars, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.clip.Source != "dev" {
		t.Errorf("expected source dev, got %s", m.clip.Source)
	}
	keys, _ := m.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestCopySelectedKeys(t *testing.T) {
	m := newManager()
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	if err := m.Copy("dev", vars, []string{"A", "C"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, _ := m.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestCopyMissingKeyError(t *testing.T) {
	m := newManager()
	vars := map[string]string{"A": "1"}
	err := m.Copy("dev", vars, []string{"MISSING"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestCopyEmptySourceError(t *testing.T) {
	m := newManager()
	err := m.Copy("", map[string]string{"A": "1"}, nil)
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestPasteOverwritesDst(t *testing.T) {
	m := newManager()
	_ = m.Copy("dev", map[string]string{"A": "new", "B": "2"}, nil)
	dst := map[string]string{"A": "old", "Z": "26"}
	out, err := m.Paste(dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %s", out["A"])
	}
	if out["Z"] != "26" {
		t.Errorf("expected Z=26, got %s", out["Z"])
	}
}

func TestPasteEmptyClipboardError(t *testing.T) {
	m := newManager()
	_, err := m.Paste(map[string]string{})
	if err != ErrClipboardEmpty {
		t.Errorf("expected ErrClipboardEmpty, got %v", err)
	}
}

func TestClearResetsClipboard(t *testing.T) {
	m := newManager()
	_ = m.Copy("dev", map[string]string{"A": "1"}, nil)
	m.Clear()
	if !m.IsEmpty() {
		t.Error("expected clipboard to be empty after Clear")
	}
}

func TestKeysSorted(t *testing.T) {
	m := newManager()
	_ = m.Copy("dev", map[string]string{"Z": "1", "A": "2", "M": "3"}, nil)
	keys, _ := m.Keys()
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestPasteDoesNotMutateDst(t *testing.T) {
	m := newManager()
	_ = m.Copy("dev", map[string]string{"X": "10"}, nil)
	dst := map[string]string{"Y": "20"}
	_, _ = m.Paste(dst)
	if _, ok := dst["X"]; ok {
		t.Error("Paste must not mutate the original dst map")
	}
}
