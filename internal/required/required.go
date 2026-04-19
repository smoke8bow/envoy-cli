package required

import "fmt"

// Violation represents a missing required key.
type Violation struct {
	Key     string
	Profile string
}

func (v Violation) Error() string {
	return fmt.Sprintf("profile %q: required key %q is missing or empty", v.Profile, v.Key)
}

// Manager checks profiles against a set of required keys.
type Manager struct {
	required map[string][]string // profile -> keys
}

func NewManager() *Manager {
	return &Manager{required: make(map[string][]string)}
}

// Set defines required keys for a profile.
func (m *Manager) Set(profile string, keys []string) {
	copy := make([]string, len(keys))
	for i, k := range keys {
		copy[i] = k
	}
	m.required[profile] = copy
}

// Get returns the required keys for a profile.
func (m *Manager) Get(profile string) []string {
	return m.required[profile]
}

// Check validates that all required keys are present and non-empty in vars.
// Returns a slice of Violations (empty means all satisfied).
func (m *Manager) Check(profile string, vars map[string]string) []Violation {
	keys, ok := m.required[profile]
	if !ok {
		return nil
	}
	var violations []Violation
	for _, k := range keys {
		v, exists := vars[k]
		if !exists || v == "" {
			violations = append(violations, Violation{Key: k, Profile: profile})
		}
	}
	return violations
}

// CheckAll validates all profiles that have requirements.
func (m *Manager) CheckAll(vars func(profile string) map[string]string) map[string][]Violation {
	result := make(map[string][]Violation)
	for profile := range m.required {
		if v := m.Check(profile, vars(profile)); len(v) > 0 {
			result[profile] = v
		}
	}
	return result
}
