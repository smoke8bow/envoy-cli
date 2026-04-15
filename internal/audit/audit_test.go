package audit

import (
	"os"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoadEmptyLog(t *testing.T) {
	dir := tempDir(t)
	log, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestRecordAndRecent(t *testing.T) {
	dir := tempDir(t)
	log, _ := Load(dir)

	if err := log.Record(EventCreate, "prod", "initial create"); err != nil {
		t.Fatalf("record failed: %v", err)
	}
	if err := log.Record(EventSwitch, "prod", ""); err != nil {
		t.Fatalf("record failed: %v", err)
	}

	if len(log.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Event != EventCreate {
		t.Errorf("expected first event to be %q, got %q", EventCreate, log.Entries[0].Event)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	log, _ := Load(dir)
	_ = log.Record(EventUpdate, "staging", "updated KEY")

	reloaded, err := Load(dir)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if len(reloaded.Entries) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(reloaded.Entries))
	}
	if reloaded.Entries[0].Profile != "staging" {
		t.Errorf("expected profile %q, got %q", "staging", reloaded.Entries[0].Profile)
	}
}

func TestRecentLimit(t *testing.T) {
	dir := tempDir(t)
	log, _ := Load(dir)

	for i := 0; i < 5; i++ {
		_ = log.Record(EventSwitch, "env", "")
	}

	recent := log.Recent(3)
	if len(recent) != 3 {
		t.Errorf("expected 3 recent entries, got %d", len(recent))
	}
}

func TestRecentAll(t *testing.T) {
	dir := tempDir(t)
	log, _ := Load(dir)
	_ = log.Record(EventDelete, "old", "")

	recent := log.Recent(10)
	if len(recent) != 1 {
		t.Errorf("expected 1 entry, got %d", len(recent))
	}
}
