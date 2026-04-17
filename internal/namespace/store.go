package namespace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// FileStore persists namespaces to a JSON file.
type FileStore struct {
	path string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{path: filepath.Join(dir, "namespaces.json")}
}

func (f *FileStore) Load() ([]Namespace, error) {
	data, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Namespace{}, nil
	}
	if err != nil {
		return nil, err
	}
	var ns []Namespace
	if err := json.Unmarshal(data, &ns); err != nil {
		return nil, err
	}
	return ns, nil
}

func (f *FileStore) Save(ns []Namespace) error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(ns, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o644)
}
