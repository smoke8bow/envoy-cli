package store_test

import (
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/store"
)

func TestSetAndGetMeta(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(filepath.Join(dir, "s.json"))
	_ = s.Set("dev", map[string]string{"A": "1"})

	if err := s.SetMeta("dev", "tags", "backend,staging"); err != nil {
		t.Fatalf("SetMeta: %v", err)
	}
	v, err := s.GetMeta("dev", "tags")
	if err != nil {
		t.Fatalf("GetMeta: %v", err)
	}
	if v != "backend,staging" {
		t.Errorf("expected 'backend,staging', got %q", v)
	}
}

func TestGetMetaMissingKey(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(filepath.Join(dir, "s.json"))
	_ = s.Set("dev", map[string]string{"A": "1"})

	v, err := s.GetMeta("dev", "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "" {
		t.Errorf("expected empty string, got %q", v)
	}
}

func TestSetMetaUnknownProfile(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(filepath.Join(dir, "s.json"))

	if err := s.SetMeta("ghost", "k", "v"); err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestDeleteMeta(t *testing.T) {
	dir := tempDir(t)
	s, _ := store.Load(filepath.Join(dir, "s.json"))
	_ = s.Set("dev", map[string]string{"A": "1"})
	_ = s.SetMeta("dev", "tags", "backend")

	if err := s.DeleteMeta("dev", "tags"); err != nil {
		t.Fatalf("DeleteMeta: %v", err)
	}
	v, _ := s.GetMeta("dev", "tags")
	if v != "" {
		t.Errorf("expected empty after delete, got %q", v)
	}
}

func TestMetaPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "s.json")
	s, _ := store.Load(path)
	_ = s.Set("dev", map[string]string{"A": "1"})
	_ = s.SetMeta("dev", "owner", "alice")

	s2, err := store.Load(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	v, _ := s2.GetMeta("dev", "owner")
	if v != "alice" {
		t.Errorf("expected 'alice' after reload, got %q", v)
	}
}
