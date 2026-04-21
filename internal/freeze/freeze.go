// Package freeze provides the ability to mark a profile as frozen,
// preventing any modifications to its environment variables.
package freeze

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ErrFrozen is returned when an operation is attempted on a frozen profile.
var ErrFrozen = errors.New("profile is frozen")

// ErrNotFrozen is returned when trying to unfreeze a profile that is not frozen.
var ErrNotFrozen = errors.New("profile is not frozen")

// Manager manages the frozen state of profiles.
type Manager struct {
	path string
	data map[string]bool
}

// NewManager creates a new Manager backed by a JSON file at the given path.
func NewManager(path string) (*Manager, error) {
	m := &Manager{path: path, data: make(map[string]bool)}
	if err := m.load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) load() error {
	b, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
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

// Freeze marks the named profile as frozen.
func (m *Manager) Freeze(profile string) error {
	if profile == "" {
		return errors.New("profile name must not be empty")
	}
	m.data[profile] = true
	return m.save()
}

// Unfreeze removes the frozen mark from the named profile.
func (m *Manager) Unfreeze(profile string) error {
	if !m.data[profile] {
		return ErrNotFrozen
	}
	delete(m.data, profile)
	return m.save()
}

// IsFrozen reports whether the named profile is currently frozen.
func (m *Manager) IsFrozen(profile string) bool {
	return m.data[profile]
}

// List returns the names of all frozen profiles.
func (m *Manager) List() []string {
	out := make([]string, 0, len(m.data))
	for k, v := range m.data {
		if v {
			out = append(out, k)
		}
	}
	return out
}

// Guard returns ErrFrozen if the profile is frozen, otherwise nil.
func (m *Manager) Guard(profile string) error {
	if m.IsFrozen(profile) {
		return ErrFrozen
	}
	return nil
}
