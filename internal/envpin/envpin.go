// Package envpin allows pinning specific environment variable keys so they
// cannot be modified or deleted during bulk operations.
package envpin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ErrNotPinned is returned when trying to unpin a key that is not pinned.
var ErrNotPinned = errors.New("key is not pinned")

// ErrAlreadyPinned is returned when trying to pin a key that is already pinned.
var ErrAlreadyPinned = errors.New("key is already pinned")

// Manager manages pinned keys per profile.
type Manager struct {
	path string
	data map[string]map[string]bool // profile -> key -> pinned
}

// NewManager creates a Manager backed by the given file path.
func NewManager(path string) (*Manager, error) {
	m := &Manager{path: path, data: make(map[string]map[string]bool)}
	if err := m.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("envpin: load: %w", err)
	}
	return m, nil
}

// Pin marks key in profile as pinned.
func (m *Manager) Pin(profile, key string) error {
	if profile == "" || key == "" {
		return errors.New("envpin: profile and key must not be empty")
	}
	if m.IsPinned(profile, key) {
		return fmt.Errorf("envpin: %q in %q: %w", key, profile, ErrAlreadyPinned)
	}
	if m.data[profile] == nil {
		m.data[profile] = make(map[string]bool)
	}
	m.data[profile][key] = true
	return m.save()
}

// Unpin removes the pin from key in profile.
func (m *Manager) Unpin(profile, key string) error {
	if !m.IsPinned(profile, key) {
		return fmt.Errorf("envpin: %q in %q: %w", key, profile, ErrNotPinned)
	}
	delete(m.data[profile], key)
	return m.save()
}

// IsPinned reports whether key is pinned in profile.
func (m *Manager) IsPinned(profile, key string) bool {
	return m.data[profile][key]
}

// Keys returns all pinned keys for profile in sorted order.
func (m *Manager) Keys(profile string) []string {
	keys := make([]string, 0, len(m.data[profile]))
	for k := range m.data[profile] {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FilterWritable removes pinned keys from vars, returning only those that may be written.
func (m *Manager) FilterWritable(profile string, vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if !m.IsPinned(profile, k) {
			out[k] = v
		}
	}
	return out
}

func (m *Manager) load() error {
	b, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &m.data)
}

func (m *Manager) save() error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, b, 0o644)
}
