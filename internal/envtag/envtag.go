// Package envtag provides key-level tagging for environment variables within a profile.
// Tags are stored as metadata alongside the profile store.
package envtag

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Store persists tag data.
type Store interface {
	GetMeta(profile, key string) (string, error)
	SetMeta(profile, key, value string) error
}

// Manager manages tags on individual env keys.
type Manager struct {
	store Store
}

// NewManager returns a Manager backed by store.
func NewManager(store Store) *Manager {
	return &Manager{store: store}
}

const metaKey = "__envtags__"

// Set assigns tags to the given key within a profile.
func (m *Manager) Set(profile, envKey string, tags []string) error {
	if profile == "" {
		return fmt.Errorf("profile name must not be empty")
	}
	if envKey == "" {
		return fmt.Errorf("env key must not be empty")
	}
	sorted := make([]string, len(tags))
	copy(sorted, tags)
	sort.Strings(sorted)
	all, err := m.load(profile)
	if err != nil {
		return err
	}
	all[envKey] = sorted
	return m.save(profile, all)
}

// Get returns the tags for the given key within a profile.
func (m *Manager) Get(profile, envKey string) ([]string, error) {
	all, err := m.load(profile)
	if err != nil {
		return nil, err
	}
	tags, ok := all[envKey]
	if !ok {
		return []string{}, nil
	}
	return tags, nil
}

// Remove clears all tags for the given key within a profile.
func (m *Manager) Remove(profile, envKey string) error {
	all, err := m.load(profile)
	if err != nil {
		return err
	}
	delete(all, envKey)
	return m.save(profile, all)
}

// List returns a map of env key -> tags for all tagged keys in a profile.
func (m *Manager) List(profile string) (map[string][]string, error) {
	return m.load(profile)
}

func (m *Manager) load(profile string) (map[string][]string, error) {
	raw, err := m.store.GetMeta(profile, metaKey)
	if err != nil || raw == "" {
		return map[string][]string{}, nil
	}
	var out map[string][]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string][]string{}, nil
	}
	return out, nil
}

func (m *Manager) save(profile string, data map[string][]string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("envtag: marshal: %w", err)
	}
	return m.store.SetMeta(profile, metaKey, string(b))
}
