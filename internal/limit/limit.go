package limit

import (
	"errors"
	"fmt"
)

// ErrLimitExceeded is returned when a profile exceeds its variable count limit.
var ErrLimitExceeded = errors.New("variable limit exceeded")

// Manager enforces maximum variable counts per profile.
type Manager struct {
	defaultMax int
	overrides  map[string]int
}

// NewManager creates a Manager with the given default max.
func NewManager(defaultMax int) *Manager {
	return &Manager{
		defaultMax: defaultMax,
		overrides:  make(map[string]int),
	}
}

// SetLimit sets a per-profile override limit.
func (m *Manager) SetLimit(profile string, max int) {
	m.overrides[profile] = max
}

// GetLimit returns the effective limit for a profile.
func (m *Manager) GetLimit(profile string) int {
	if v, ok := m.overrides[profile]; ok {
		return v
	}
	return m.defaultMax
}

// Check returns an error if the given variable count exceeds the profile limit.
func (m *Manager) Check(profile string, count int) error {
	lim := m.GetLimit(profile)
	if count > lim {
		return fmt.Errorf("%w: profile %q has %d variables, max is %d",
			ErrLimitExceeded, profile, count, lim)
	}
	return nil
}

// CheckVars is a convenience wrapper that checks len(vars) against the limit.
func (m *Manager) CheckVars(profile string, vars map[string]string) error {
	return m.Check(profile, len(vars))
}
