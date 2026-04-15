package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/store"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envoy-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoadEmptyStore(t *testing.T) {
	dir := tempDir(t)
	s, err := store.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Errorf("expected empty store, got %d sets", len(s.List()))
	}
}

func TestAddAndGet(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(dir)

	set := store.EnvSet{Name: "dev", Vars: map[string]string{"FOO": "bar", "PORT": "8080"}}
	s.Add(set)

	got, ok := s.Get("dev")
	if !ok {
		t.Fatal("expected to find 'dev' set")
	}
	if got.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", got.Vars["FOO"])
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(dir)
	s.Add(store.EnvSet{Name: "prod", Vars: map[string]string{"ENV": "production"}})

	if err := s.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".envoy.json")); err != nil {
		t.Fatalf("store file not created: %v", err)
	}

	s2, err := store.Load(dir)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	got, ok := s2.Get("prod")
	if !ok {
		t.Fatal("expected 'prod' after reload")
	}
	if got.Vars["ENV"] != "production" {
		t.Errorf("expected ENV=production, got %s", got.Vars["ENV"])
	}
}

func TestDelete(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(dir)
	s.Add(store.EnvSet{Name: "staging", Vars: map[string]string{}})

	if !s.Delete("staging") {
		t.Error("expected delete to return true")
	}
	if s.Delete("staging") {
		t.Error("expected delete of missing key to return false")
	}
	if _, ok := s.Get("staging"); ok {
		t.Error("expected 'staging' to be gone")
	}
}
