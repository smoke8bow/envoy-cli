package envvault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/envvault"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envvault-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestFileStorageLoadEmpty(t *testing.T) {
	dir := tempDir(t)
	fs := envvault.NewFileStorage(filepath.Join(dir, "vault.json"))
	blobs, err := fs.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(blobs) != 0 {
		t.Errorf("expected empty map, got %v", blobs)
	}
}

func TestFileStorageSaveAndLoad(t *testing.T) {
	dir := tempDir(t)
	fs := envvault.NewFileStorage(filepath.Join(dir, "vault.json"))
	input := map[string]string{"prod": "encryptedblob=="}
	if err := fs.Save(input); err != nil {
		t.Fatalf("Save: %v", err)
	}
	blobs, err := fs.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if blobs["prod"] != "encryptedblob==" {
		t.Errorf("unexpected blob: %q", blobs["prod"])
	}
}

func TestFileStorageCreatesParentDir(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "nested", "deep", "vault.json")
	fs := envvault.NewFileStorage(path)
	if err := fs.Save(map[string]string{"k": "v"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestFileStorageRoundtripWithManager(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "vault.json")
	fs := envvault.NewFileStorage(path)
	m := envvault.NewManager(fs, "mypassphrase")
	vars := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
	if err := m.Put("staging", vars); err != nil {
		t.Fatalf("Put: %v", err)
	}
	// Reload manager from same file to verify persistence.
	m2 := envvault.NewManager(envvault.NewFileStorage(path), "mypassphrase")
	got, err := m2.Get("staging")
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if got["DB_PASS"] != "s3cr3t" {
		t.Errorf("DB_PASS mismatch: %q", got["DB_PASS"])
	}
}
