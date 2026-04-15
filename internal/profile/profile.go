package profile

import (
	"errors"
	"fmt"

	"envoy-cli/internal/store"
)

// ErrProfileNotFound is returned when a named profile does not exist.
var ErrProfileNotFound = errors.New("profile not found")

// ErrProfileExists is returned when attempting to create a duplicate profile.
var ErrProfileExists = errors.New("profile already exists")

// Manager handles CRUD operations for named environment profiles.
type Manager struct {
	store *store.Store
}

// NewManager creates a Manager backed by the given Store.
func NewManager(s *store.Store) *Manager {
	return &Manager{store: s}
}

// Create adds a new named profile with the provided variables.
// Returns ErrProfileExists if a profile with that name already exists.
func (m *Manager) Create(name string, vars map[string]string) error {
	if _, err := m.store.Get(name); err == nil {
		return fmt.Errorf("%w: %s", ErrProfileExists, name)
	}
	return m.store.Add(name, vars)
}

// Get retrieves the variables for a named profile.
// Returns ErrProfileNotFound if the profile does not exist.
func (m *Manager) Get(name string) (map[string]string, error) {
	vars, err := m.store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	return vars, nil
}

// Update replaces the variables for an existing profile.
// Returns ErrProfileNotFound if the profile does not exist.
func (m *Manager) Update(name string, vars map[string]string) error {
	if _, err := m.store.Get(name); err != nil {
		return fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	if err := m.store.Delete(name); err != nil {
		return err
	}
	return m.store.Add(name, vars)
}

// Delete removes a named profile.
// Returns ErrProfileNotFound if the profile does not exist.
func (m *Manager) Delete(name string) error {
	if _, err := m.store.Get(name); err != nil {
		return fmt.Errorf("%w: %s", ErrProfileNotFound, name)
	}
	return m.store.Delete(name)
}

// List returns all profile names currently stored.
func (m *Manager) List() []string {
	return m.store.List()
}
