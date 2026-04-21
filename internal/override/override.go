// Package override provides a mechanism to layer ephemeral key-value pairs
// on top of a profile without persisting the changes to the store.
package override

import "fmt"

// Layer holds a named set of ephemeral overrides for a profile.
type Layer struct {
	Profile string
	Vars    map[string]string
}

// Manager manages override layers keyed by profile name.
type Manager struct {
	layers map[string]map[string]string
}

// NewManager returns an empty Manager.
func NewManager() *Manager {
	return &Manager{layers: make(map[string]map[string]string)}
}

// Set adds or replaces a key in the override layer for the given profile.
func (m *Manager) Set(profile, key, value string) error {
	if profile == "" {
		return fmt.Errorf("override: profile name must not be empty")
	}
	if key == "" {
		return fmt.Errorf("override: key must not be empty")
	}
	if _, ok := m.layers[profile]; !ok {
		m.layers[profile] = make(map[string]string)
	}
	m.layers[profile][key] = value
	return nil
}

// Unset removes a key from the override layer for the given profile.
// Returns an error if the key does not exist in the layer.
func (m *Manager) Unset(profile, key string) error {
	if profile == "" {
		return fmt.Errorf("override: profile name must not be empty")
	}
	layer, ok := m.layers[profile]
	if !ok {
		return fmt.Errorf("override: no layer for profile %q", profile)
	}
	if _, ok := layer[key]; !ok {
		return fmt.Errorf("override: key %q not found in layer for profile %q", key, profile)
	}
	delete(layer, key)
	return nil
}

// Apply merges the override layer for the given profile on top of base.
// base is not mutated; a new map is returned.
func (m *Manager) Apply(profile string, base map[string]string) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range m.layers[profile] {
		result[k] = v
	}
	return result
}

// Clear removes the entire override layer for the given profile.
func (m *Manager) Clear(profile string) {
	delete(m.layers, profile)
}

// Layer returns a copy of the current override layer for the given profile.
func (m *Manager) Layer(profile string) map[string]string {
	copy := make(map[string]string)
	for k, v := range m.layers[profile] {
		copy[k] = v
	}
	return copy
}
