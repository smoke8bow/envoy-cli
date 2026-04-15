package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const storeFileName = ".envoy.json"

// EnvSet represents a named set of environment variables.
type EnvSet struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Store holds all named environment sets for a project.
type Store struct {
	Sets map[string]EnvSet `json:"sets"`
	path string
}

// Load reads the store from the given directory, or returns an empty store.
func Load(dir string) (*Store, error) {
	path := filepath.Join(dir, storeFileName)
	s := &Store{
		Sets: make(map[string]EnvSet),
		path: path,
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

// Save persists the store to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

// Add inserts or replaces a named env set.
func (s *Store) Add(set EnvSet) {
	s.Sets[set.Name] = set
}

// Get returns a named env set, or false if not found.
func (s *Store) Get(name string) (EnvSet, bool) {
	set, ok := s.Sets[name]
	return set, ok
}

// Delete removes a named env set.
func (s *Store) Delete(name string) bool {
	_, ok := s.Sets[name]
	delete(s.Sets, name)
	return ok
}

// List returns all env set names.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.Sets))
	for name := range s.Sets {
		names = append(names, name)
	}
	return names
}
