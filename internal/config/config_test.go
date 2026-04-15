package config

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envoy-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoadDefaults(t *testing.T) {
	dir := tempDir(t)
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := DefaultConfig()
	if cfg.AuditEnabled != def.AuditEnabled {
		t.Errorf("AuditEnabled: got %v, want %v", cfg.AuditEnabled, def.AuditEnabled)
	}
	if cfg.EncryptionEnabled != def.EncryptionEnabled {
		t.Errorf("EncryptionEnabled: got %v, want %v", cfg.EncryptionEnabled, def.EncryptionEnabled)
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := tempDir(t)
	cfg := Config{
		DefaultShell:      "zsh",
		StorePath:         "/tmp/store",
		AuditEnabled:      false,
		EncryptionEnabled: true,
		AutoExport:        true,
	}
	if err := Save(dir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.DefaultShell != cfg.DefaultShell {
		t.Errorf("DefaultShell: got %q, want %q", loaded.DefaultShell, cfg.DefaultShell)
	}
	if loaded.EncryptionEnabled != cfg.EncryptionEnabled {
		t.Errorf("EncryptionEnabled: got %v, want %v", loaded.EncryptionEnabled, cfg.EncryptionEnabled)
	}
	if loaded.AutoExport != cfg.AutoExport {
		t.Errorf("AutoExport: got %v, want %v", loaded.AutoExport, cfg.AutoExport)
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := filepath.Join(tempDir(t), "nested", "dir")
	if err := Save(dir, DefaultConfig()); err != nil {
		t.Fatalf("Save should create nested dirs: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, configFileName)); err != nil {
		t.Errorf("config file not found after Save: %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, []byte("not-json{"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	_, err := Load(dir)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
