package store

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// Store holds named environment variable sets and persists them to disk.
type Store struct {
	path string
	data map[string]map[string]string
}

// Load reads the store from path, or returns an empty store if the file does
// not yet exist.
func Load(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]map[string]string),
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(bytes, &s.data); err != nil {
		return nil, err
	}
	return s, nil
}

// save persists the current state to disk.
func (s *Store) save() error {
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, bytes, 0o600)
}

// Add inserts a new named profile. It does not check for duplicates.
func (s *Store) Add(name string, vars map[string]string) error {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	s.data[name] = copy
	return s.save()
}

// Get returns the variables for the named profile, or an error if not found.
func (s *Store) Get(name string) (map[string]string, error) {
	vars, ok := s.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return copy, nil
}

// Delete removes the named profile from the store.
func (s *Store) Delete(name string) error {
	delete(s.data, name)
	return s.save()
}

// List returns all profile names in sorted order.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.data))
	for name := range s.data {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
