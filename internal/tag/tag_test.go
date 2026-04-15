package tag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/store"
	"github.com/envoy-cli/envoy-cli/internal/tag"
)

func newManager(t *testing.T) *tag.Manager {
	t.Helper()
	dir, err := os.MkdirTemp("", "tag-test-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := store.Load(filepath.Join(dir, "store.json"))
	if err != nil {
		t.Fatalf("store.Load: %v", err)
	}
	// seed a profile so SetMeta / GetMeta have something to work with
	_ = s.Set("dev", map[string]string{"X": "1"})
	_ = s.Set("prod", map[string]string{"X": "2"})
	return tag.NewManager(s)
}

func TestAddAndList(t *testing.T) {
	m := newManager(t)
	if err := m.Add("dev", "backend"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	tags, err := m.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tags) != 1 || tags[0] != "backend" {
		t.Errorf("expected [backend], got %v", tags)
	}
}

func TestAddDuplicateIgnored(t *testing.T) {
	m := newManager(t)
	_ = m.Add("dev", "backend")
	_ = m.Add("dev", "backend")
	tags, _ := m.List("dev")
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestRemoveTag(t *testing.T) {
	m := newManager(t)
	_ = m.Add("dev", "backend")
	_ = m.Add("dev", "staging")
	if err := m.Remove("dev", "staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	tags, _ := m.List("dev")
	if len(tags) != 1 || tags[0] != "backend" {
		t.Errorf("expected [backend], got %v", tags)
	}
}

func TestRemoveNotFound(t *testing.T) {
	m := newManager(t)
	err := m.Remove("dev", "ghost")
	if err == nil {
		t.Fatal("expected error removing missing tag")
	}
}

func TestProfilesWithTag(t *testing.T) {
	m := newManager(t)
	_ = m.Add("dev", "backend")
	_ = m.Add("prod", "backend")
	_ = m.Add("prod", "live")

	profiles, err := m.ProfilesWithTag("backend")
	if err != nil {
		t.Fatalf("ProfilesWithTag: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d: %v", len(profiles), profiles)
	}
}

func TestListEmpty(t *testing.T) {
	m := newManager(t)
	tags, err := m.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty list, got %v", tags)
	}
}
