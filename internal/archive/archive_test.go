package archive

import (
	"os"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "archive-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestListEmpty(t *testing.T) {
	mgr := NewManager(tempDir(t))
	entries, err := mgr.List("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestArchiveAndList(t *testing.T) {
	mgr := NewManager(tempDir(t))
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := mgr.Archive("dev", vars); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	entries, err := mgr.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Profile != "dev" {
		t.Errorf("expected profile dev, got %s", entries[0].Profile)
	}
	if entries[0].Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestLatestReturnsNewest(t *testing.T) {
	mgr := NewManager(tempDir(t))
	if err := mgr.Archive("prod", map[string]string{"V": "1"}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := mgr.Archive("prod", map[string]string{"V": "2"}); err != nil {
		t.Fatal(err)
	}
	entry, err := mgr.Latest("prod")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if entry.Vars["V"] != "2" {
		t.Errorf("expected V=2, got %s", entry.Vars["V"])
	}
}

func TestLatestNotFound(t *testing.T) {
	mgr := NewManager(tempDir(t))
	_, err := mgr.Latest("ghost")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestListAllProfiles(t *testing.T) {
	mgr := NewManager(tempDir(t))
	_ = mgr.Archive("dev", map[string]string{"A": "1"})
	_ = mgr.Archive("prod", map[string]string{"B": "2"})
	all, err := mgr.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
