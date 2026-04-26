package envlabel_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/envlabel"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "envlabel-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func newManager(t *testing.T) *envlabel.Manager {
	t.Helper()
	m, err := envlabel.NewManager(filepath.Join(tempDir(t), "labels.json"))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestSetAndGet(t *testing.T) {
	m := newManager(t)
	if err := m.Set("prod", "team", "backend"); err != nil {
		t.Fatal(err)
	}
	if err := m.Set("prod", "env", "production"); err != nil {
		t.Fatal(err)
	}
	labels := m.Get("prod")
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	// sorted by key: env, team
	if labels[0].Key != "env" || labels[0].Value != "production" {
		t.Errorf("unexpected first label: %+v", labels[0])
	}
	if labels[1].Key != "team" || labels[1].Value != "backend" {
		t.Errorf("unexpected second label: %+v", labels[1])
	}
}

func TestGetUnknownProfile(t *testing.T) {
	m := newManager(t)
	if got := m.Get("ghost"); len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestSetEmptyProfileError(t *testing.T) {
	m := newManager(t)
	if err := m.Set("", "k", "v"); err == nil {
		t.Error("expected error for empty profile")
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	m := newManager(t)
	if err := m.Set("prod", "", "v"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestRemoveLabel(t *testing.T) {
	m := newManager(t)
	_ = m.Set("prod", "team", "backend")
	_ = m.Set("prod", "env", "staging")
	if err := m.Remove("prod", "team"); err != nil {
		t.Fatal(err)
	}
	labels := m.Get("prod")
	if len(labels) != 1 || labels[0].Key != "env" {
		t.Errorf("unexpected labels after remove: %v", labels)
	}
}

func TestRemoveNotFound(t *testing.T) {
	m := newManager(t)
	if err := m.Remove("prod", "missing"); err == nil {
		t.Error("expected error removing non-existent label")
	}
}

func TestPersistence(t *testing.T) {
	path := filepath.Join(tempDir(t), "labels.json")
	m1, _ := envlabel.NewManager(path)
	_ = m1.Set("dev", "owner", "alice")

	m2, err := envlabel.NewManager(path)
	if err != nil {
		t.Fatal(err)
	}
	labels := m2.Get("dev")
	if len(labels) != 1 || labels[0].Value != "alice" {
		t.Errorf("label not persisted: %v", labels)
	}
}

func TestProfiles(t *testing.T) {
	m := newManager(t)
	_ = m.Set("prod", "env", "production")
	_ = m.Set("dev", "env", "development")
	profiles := m.Profiles()
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
	if profiles[0] != "dev" || profiles[1] != "prod" {
		t.Errorf("unexpected order: %v", profiles)
	}
}
