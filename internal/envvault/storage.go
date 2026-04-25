package envvault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FileStorage persists encrypted blobs as a JSON file on disk.
type FileStorage struct {
	path string
}

// NewFileStorage returns a FileStorage that reads/writes to path.
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

// Load reads the blob map from disk. Returns an empty map if the file does not exist.
func (f *FileStorage) Load() (map[string]string, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("envvault storage: read %s: %w", f.path, err)
	}
	var blobs map[string]string
	if err := json.Unmarshal(data, &blobs); err != nil {
		return nil, fmt.Errorf("envvault storage: unmarshal: %w", err)
	}
	return blobs, nil
}

// Save writes the blob map to disk as JSON, creating parent directories as needed.
func (f *FileStorage) Save(blobs map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o700); err != nil {
		return fmt.Errorf("envvault storage: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(blobs, "", "  ")
	if err != nil {
		return fmt.Errorf("envvault storage: marshal: %w", err)
	}
	if err := os.WriteFile(f.path, data, 0o600); err != nil {
		return fmt.Errorf("envvault storage: write %s: %w", f.path, err)
	}
	return nil
}
