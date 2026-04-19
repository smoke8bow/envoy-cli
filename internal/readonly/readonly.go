package readonly

import "fmt"

// Manager tracks which profiles are marked read-only.
type Manager struct {
	flags map[string]bool
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	return &Manager{flags: make(map[string]bool)}
}

// Set marks a profile as read-only.
func (m *Manager) Set(profile string) error {
	if profile == "" {
		return fmt.Errorf("readonly: profile name must not be empty")
	}
	m.flags[profile] = true
	return nil
}

// Unset removes the read-only flag from a profile.
func (m *Manager) Unset(profile string) error {
	if !m.flags[profile] {
		return fmt.Errorf("readonly: profile %q is not read-only", profile)
	}
	delete(m.flags, profile)
	return nil
}

// IsReadOnly reports whether the given profile is marked read-only.
func (m *Manager) IsReadOnly(profile string) bool {
	return m.flags[profile]
}

// Check returns an error if the profile is read-only.
func (m *Manager) Check(profile string) error {
	if m.flags[profile] {
		return fmt.Errorf("readonly: profile %q is read-only and cannot be modified", profile)
	}
	return nil
}

// List returns all profiles currently marked read-only.
func (m *Manager) List() []string {
	out := make([]string, 0, len(m.flags))
	for k := range m.flags {
		out = append(out, k)
	}
	return out
}
