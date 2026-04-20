// Package defaults manages per-profile default key-value pairs that are
// automatically merged into a profile when it is loaded or applied.
package defaults

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const filename = "defaults.json"

// Manager handles storage and retrieval of default env vars per profile.
type Manager struct {
	path string
}

// entry maps profile name → default key/value pairs.
type entry map[string]map[string]string

// NewManager creates a Manager whose state is persisted in dir.
func NewManager(dir string) *Manager {
	return &Manager{path: filepath.Join(dir, filename)}
}

func (m *Manager) load() (entry, error) {
	data, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return make(entry), nil
	}
	if err != nil {
		return nil, err
	}
	var e entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return e, nil
}

func (m *Manager) save(e entry) error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o644)
}

// Set stores defaults for the given profile, replacing any existing defaults.
func (m *Manager) Set(profile string, vars map[string]string) error {
	if profile == "" {
		return errors.New("defaults: profile name must not be empty")
	}
	e, err := m.load()
	if err != nil {
		return err
	}
	e[profile] = vars
	return m.save(e)
}

// Get returns the defaults registered for the given profile.
// Returns an empty map (not an error) when no defaults are set.
func (m *Manager) Get(profile string) (map[string]string, error) {
	e, err := m.load()
	if err != nil {
		return nil, err
	}
	if v, ok := e[profile]; ok {
		return v, nil
	}
	return map[string]string{}, nil
}

// Apply merges the registered defaults for profile into vars.
// Keys already present in vars are NOT overwritten.
func (m *Manager) Apply(profile string, vars map[string]string) (map[string]string, error) {
	defaults, err := m.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("defaults: %w", err)
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	for k, v := range defaults {
		if _, exists := out[k]; !exists {
			out[k] = v
		}
	}
	return out, nil
}

// Delete removes all defaults for the given profile.
func (m *Manager) Delete(profile string) error {
	e, err := m.load()
	if err != nil {
		return err
	}
	delete(e, profile)
	return m.save(e)
}
