package alias

import (
	"errors"
	"regexp"
)

var validAlias = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Store is the minimal interface needed by the Manager.
type Store interface {
	List() ([]string, error)
}

// Manager manages short aliases that map to profile names.
type Manager struct {
	aliases map[string]string
	store   Store
}

// NewManager creates a new alias Manager.
func NewManager(s Store) *Manager {
	return &Manager{
		aliases: make(map[string]string),
		store:   s,
	}
}

// Set registers an alias pointing to a profile name.
func (m *Manager) Set(alias, profile string) error {
	if alias == "" {
		return errors.New("alias: name must not be empty")
	}
	if !validAlias.MatchString(alias) {
		return errors.New("alias: name contains invalid characters")
	}
	if profile == "" {
		return errors.New("alias: profile must not be empty")
	}
	profiles, err := m.store.List()
	if err != nil {
		return err
	}
	found := false
	for _, p := range profiles {
		if p == profile {
			found = true
			break
		}
	}
	if !found {
		return errors.New("alias: profile not found: " + profile)
	}
	m.aliases[alias] = profile
	return nil
}

// Resolve returns the profile name for the given alias.
func (m *Manager) Resolve(alias string) (string, error) {
	p, ok := m.aliases[alias]
	if !ok {
		return "", errors.New("alias: not found: " + alias)
	}
	return p, nil
}

// Remove deletes an alias.
func (m *Manager) Remove(alias string) error {
	if _, ok := m.aliases[alias]; !ok {
		return errors.New("alias: not found: " + alias)
	}
	delete(m.aliases, alias)
	return nil
}

// List returns all registered aliases.
func (m *Manager) List() map[string]string {
	out := make(map[string]string, len(m.aliases))
	for k, v := range m.aliases {
		out[k] = v
	}
	return out
}
