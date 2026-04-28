// Package envaliases provides key-alias mapping for environment variable profiles.
// It allows defining short aliases for long environment variable keys.
package envaliases

import (
	"errors"
	"fmt"
	"sort"
)

// AliasMap maps alias -> canonical key.
type AliasMap map[string]string

// Manager manages aliases for a named profile.
type Manager struct {
	aliases map[string]AliasMap // profile -> alias map
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	return &Manager{aliases: make(map[string]AliasMap)}
}

// Set registers an alias for a canonical key within a profile.
// Returns an error if profile or alias is empty.
func (m *Manager) Set(profile, alias, key string) error {
	if profile == "" {
		return errors.New("profile name must not be empty")
	}
	if alias == "" {
		return errors.New("alias must not be empty")
	}
	if key == "" {
		return errors.New("canonical key must not be empty")
	}
	if _, ok := m.aliases[profile]; !ok {
		m.aliases[profile] = make(AliasMap)
	}
	m.aliases[profile][alias] = key
	return nil
}

// Resolve returns the canonical key for the given alias in a profile.
// Returns an error if the alias is not registered.
func (m *Manager) Resolve(profile, alias string) (string, error) {
	am, ok := m.aliases[profile]
	if !ok {
		return "", fmt.Errorf("no aliases registered for profile %q", profile)
	}
	key, ok := am[alias]
	if !ok {
		return "", fmt.Errorf("alias %q not found in profile %q", alias, profile)
	}
	return key, nil
}

// List returns all aliases registered for a profile, sorted.
func (m *Manager) List(profile string) []string {
	am := m.aliases[profile]
	aliases := make([]string, 0, len(am))
	for a := range am {
		aliases = append(aliases, a)
	}
	sort.Strings(aliases)
	return aliases
}

// Remove deletes an alias from a profile. Returns an error if not found.
func (m *Manager) Remove(profile, alias string) error {
	am, ok := m.aliases[profile]
	if !ok {
		return fmt.Errorf("no aliases registered for profile %q", profile)
	}
	if _, ok := am[alias]; !ok {
		return fmt.Errorf("alias %q not found in profile %q", alias, profile)
	}
	delete(am, alias)
	return nil
}

// Expand replaces any keys in vars that match registered aliases with their
// canonical equivalents. Non-matching keys are left unchanged.
func (m *Manager) Expand(profile string, vars map[string]string) map[string]string {
	am := m.aliases[profile]
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if canonical, ok := am[k]; ok {
			out[canonical] = v
		} else {
			out[k] = v
		}
	}
	return out
}
