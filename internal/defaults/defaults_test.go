package defaults

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "defaults-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newManager(t *testing.T) *Manager {
	t.Helper()
	return NewManager(tempDir(t))
}

func TestGetEmptyReturnsEmpty(t *testing.T) {
	m := newManager(t)
	vars, err := m.Get("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vars) != 0 {
		t.Errorf("expected empty map, got %v", vars)
	}
}

func TestSetAndGet(t *testing.T) {
	m := newManager(t)
	input := map[string]string{"LOG_LEVEL": "info", "TIMEOUT": "30"}
	if err := m.Set("dev", input); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	for k, v := range input {
		if got[k] != v {
			t.Errorf("key %q: want %q, got %q", k, v, got[k])
		}
	}
}

func TestSetEmptyProfileError(t *testing.T) {
	m := newManager(t)
	if err := m.Set("", map[string]string{"K": "V"}); err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestApplyDoesNotOverwriteExisting(t *testing.T) {
	m := newManager(t)
	_ = m.Set("dev", map[string]string{"LOG_LEVEL": "debug", "REGION": "us-east-1"})

	vars := map[string]string{"LOG_LEVEL": "warn"}
	out, err := m.Apply("dev", vars)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out["LOG_LEVEL"] != "warn" {
		t.Errorf("existing key overwritten: got %q", out["LOG_LEVEL"])
	}
	if out["REGION"] != "us-east-1" {
		t.Errorf("default not injected: got %q", out["REGION"])
	}
}

func TestApplyNoDefaults(t *testing.T) {
	m := newManager(t)
	vars := map[string]string{"A": "1"}
	out, err := m.Apply("dev", vars)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if out["A"] != "1" || len(out) != 1 {
		t.Errorf("unexpected result: %v", out)
	}
}

func TestDeleteRemovesDefaults(t *testing.T) {
	m := newManager(t)
	_ = m.Set("dev", map[string]string{"X": "1"})
	if err := m.Delete("dev"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, _ := m.Get("dev")
	if len(got) != 0 {
		t.Errorf("expected empty after delete, got %v", got)
	}
}

func TestPerscrossManagers(t *testing.T) {
	dir := tempDir(t)
	m1 := NewManager(dir)
	_ = m1.Set("prod"ENV": "production"})

	m2 := NewManager(dir)
	got, err := m2.Get("prod")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got["ENV"] != "production" {
		t.Errorf("persistence failed: got %v", got)
	}
}

func TestStateLocation(t *testing.T) {
	dir := tempDir(t)
	m := NewManager(dir)
	_ = m.Set("x", map[string]string{"K": "V"})
	expected := filepath.Join(dir, "defaults.json")
	if _, err := os.Stat(expected); err != nil {
		t.Errorf("expected file at %s: %v", expected, err)
	}
}
