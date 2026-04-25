// Package envflag provides boolean feature-flag management backed by
// environment variable profiles. Flags are stored as "true"/"false" string
// values inside a named env map.
package envflag

import (
	"errors"
	"fmt"
	"strings"
)

// Getter retrieves a named profile's environment variables.
type Getter interface {
	Get(name string) (map[string]string, error)
}

// Setter persists a named profile's environment variables.
type Setter interface {
	Set(name string, vars map[string]string) error
}

// Store combines Getter and Setter.
type Store interface {
	Getter
	Setter
}

// Manager manages feature flags within an env profile.
type Manager struct {
	store   Store
	profile string
}

// NewManager returns a Manager that stores flags in the given profile.
func NewManager(store Store, profile string) (*Manager, error) {
	if profile == "" {
		return nil, errors.New("envflag: profile name must not be empty")
	}
	return &Manager{store: store, profile: profile}, nil
}

// Set sets the named flag to enabled or disabled.
func (m *Manager) Set(flag string, enabled bool) error {
	if err := validateFlag(flag); err != nil {
		return err
	}
	vars, err := m.load()
	if err != nil {
		return err
	}
	if enabled {
		vars[flag] = "true"
	} else {
		vars[flag] = "false"
	}
	return m.store.Set(m.profile, vars)
}

// IsEnabled reports whether the named flag is enabled.
// Returns false if the flag does not exist.
func (m *Manager) IsEnabled(flag string) (bool, error) {
	if err := validateFlag(flag); err != nil {
		return false, err
	}
	vars, err := m.load()
	if err != nil {
		return false, err
	}
	v, ok := vars[flag]
	if !ok {
		return false, nil
	}
	return strings.EqualFold(v, "true"), nil
}

// Delete removes a flag from the profile.
func (m *Manager) Delete(flag string) error {
	if err := validateFlag(flag); err != nil {
		return err
	}
	vars, err := m.load()
	if err != nil {
		return err
	}
	if _, ok := vars[flag]; !ok {
		return fmt.Errorf("envflag: flag %q not found", flag)
	}
	delete(vars, flag)
	return m.store.Set(m.profile, vars)
}

// List returns all flags and their current boolean state.
func (m *Manager) List() (map[string]bool, error) {
	vars, err := m.load()
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(vars))
	for k, v := range vars {
		out[k] = strings.EqualFold(v, "true")
	}
	return out, nil
}

func (m *Manager) load() (map[string]string, error) {
	vars, err := m.store.Get(m.profile)
	if err != nil {
		return map[string]string{}, nil
	}
	cloned := make(map[string]string, len(vars))
	for k, v := range vars {
		cloned[k] = v
	}
	return cloned, nil
}

func validateFlag(flag string) error {
	if strings.TrimSpace(flag) == "" {
		return errors.New("envflag: flag name must not be empty")
	}
	return nil
}
