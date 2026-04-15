package switch_

import (
	"path/filepath"
	"testing"
)

func TestLoadHistoryEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	h, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if h.Last() != nil {
		t.Error("expected nil Last on empty history")
	}
	if len(h.Entries()) != 0 {
		t.Error("expected empty entries")
	}
}

func TestRecordAndLast(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	h, _ := LoadHistory(path)

	if err := h.Record("dev"); err != nil {
		t.Fatalf("Record dev: %v", err)
	}
	if err := h.Record("prod"); err != nil {
		t.Fatalf("Record prod: %v", err)
	}

	last := h.Last()
	if last == nil {
		t.Fatal("expected Last to be non-nil")
	}
	if last.Profile != "prod" {
		t.Errorf("expected last profile=prod, got %q", last.Profile)
	}
}

func TestHistoryPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	h, _ := LoadHistory(path)
	_ = h.Record("staging")
	_ = h.Record("dev")

	// Reload from disk
	h2, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("reload LoadHistory: %v", err)
	}
	if len(h2.Entries()) != 2 {
		t.Errorf("expected 2 entries after reload, got %d", len(h2.Entries()))
	}
	if h2.Last().Profile != "dev" {
		t.Errorf("expected last=dev after reload, got %q", h2.Last().Profile)
	}
}

func TestEntriesOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	h, _ := LoadHistory(path)
	profiles := []string{"a", "b", "c"}
	for _, p := range profiles {
		_ = h.Record(p)
	}
	entries := h.Entries()
	for i, p := range profiles {
		if entries[i].Profile != p {
			t.Errorf("entry[%d]: expected %q, got %q", i, p, entries[i].Profile)
		}
	}
}
