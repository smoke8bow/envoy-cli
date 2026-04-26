// Package envlabel provides key-value label management for profiles.
// Labels are arbitrary metadata (e.g. team=backend, env=staging) attached
// to a profile name and stored separately from the env vars themselves.
package envlabel

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Label is a single key=value pair attached to a profile.
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Manager manages labels for profiles.
type Manager struct {
	path string
	data map[string]map[string]string // profile -> key -> value
}

// NewManager creates a Manager backed by the given file path.
func NewManager(path string) (*Manager, error) {
	m := &Manager{path: path, data: make(map[string]map[string]string)}
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
		return fmt.Errorf("envlabel: read %s: %w", m.path, err)
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

// Set attaches a label key=value to the given profile.
func (m *Manager) Set(profile, key, value string) error {
	if profile == "" {
		return errors.New("envlabel: profile name must not be empty")
	}
	if key == "" {
		return errors.New("envlabel: label key must not be empty")
	}
	if m.data[profile] == nil {
		m.data[profile] = make(map[string]string)
	}
	m.data[profile][key] = value
	return m.save()
}

// Get returns all labels for a profile, sorted by key.
func (m *Manager) Get(profile string) []Label {
	kv := m.data[profile]
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]Label, 0, len(keys))
	for _, k := range keys {
		out = append(out, Label{Key: k, Value: kv[k]})
	}
	return out
}

// Remove deletes a single label from a profile.
func (m *Manager) Remove(profile, key string) error {
	if m.data[profile] == nil {
		return fmt.Errorf("envlabel: no labels for profile %q", profile)
	}
	if _, ok := m.data[profile][key]; !ok {
		return fmt.Errorf("envlabel: label %q not found on profile %q", key, profile)
	}
	delete(m.data[profile], key)
	if len(m.data[profile]) == 0 {
		delete(m.data, profile)
	}
	return m.save()
}

// Profiles returns all profile names that carry at least one label.
func (m *Manager) Profiles() []string {
	out := make([]string, 0, len(m.data))
	for p := range m.data {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}
