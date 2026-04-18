package group

import (
	"errors"
	"sort"
)

// Store is the minimal interface required by Manager.
type Store interface {
	List() ([]string, error)
}

// Manager manages named groups of profiles.
type Manager struct {
	store  Store
	groups map[string][]string
}

// NewManager returns a Manager backed by the given store.
func NewManager(store Store) *Manager {
	return &Manager{
		store:  store,
		groups: make(map[string][]string),
	}
}

// Create creates a new empty group.
func (m *Manager) Create(name string) error {
	if name == "" {
		return errors.New("group name must not be empty")
	}
	if _, ok := m.groups[name]; ok {
		return errors.New("group already exists: " + name)
	}
	m.groups[name] = []string{}
	return nil
}

// Add adds a profile to a group.
func (m *Manager) Add(group, profile string) error {
	profiles, ok := m.groups[group]
	if !ok {
		return errors.New("group not found: " + group)
	}
	for _, p := range profiles {
		if p == profile {
			return nil // already present
		}
	}
	m.groups[group] = append(profiles, profile)
	return nil
}

// Remove removes a profile from a group.
func (m *Manager) Remove(group, profile string) error {
	profiles, ok := m.groups[group]
	if !ok {
		return errors.New("group not found: " + group)
	}
	filtered := profiles[:0]
	for _, p := range profiles {
		if p != profile {
			filtered = append(filtered, p)
		}
	}
	m.groups[group] = filtered
	return nil
}

// Members returns sorted members of a group.
func (m *Manager) Members(group string) ([]string, error) {
	profiles, ok := m.groups[group]
	if !ok {
		return nil, errors.New("group not found: " + group)
	}
	out := make([]string, len(profiles))
	copy(out, profiles)
	sort.Strings(out)
	return out, nil
}

// List returns all group names sorted.
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.groups))
	for k := range m.groups {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Delete removes a group entirely.
func (m *Manager) Delete(name string) error {
	if _, ok := m.groups[name]; !ok {
		return errors.New("group not found: " + name)
	}
	delete(m.groups, name)
	return nil
}
