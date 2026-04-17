package quota

import (
	"errors"
	"fmt"
)

// ErrQuotaExceeded is returned when a profile exceeds its allowed variable count.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Policy defines quota rules for profiles.
type Policy struct {
	DefaultMax int
	Overrides  map[string]int
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		DefaultMax: 50,
		Overrides:  make(map[string]int),
	}
}

// Manager enforces variable count quotas per profile.
type Manager struct {
	policy Policy
}

// NewManager creates a Manager with the given policy.
func NewManager(p Policy) *Manager {
	return &Manager{policy: p}
}

// MaxFor returns the maximum number of variables allowed for a profile.
func (m *Manager) MaxFor(profile string) int {
	if v, ok := m.policy.Overrides[profile]; ok {
		return v
	}
	return m.policy.DefaultMax
}

// Check returns an error if vars exceeds the quota for the given profile.
func (m *Manager) Check(profile string, vars map[string]string) error {
	max := m.MaxFor(profile)
	if len(vars) > max {
		return fmt.Errorf("%w: profile %q has %d variables, limit is %d",
			ErrQuotaExceeded, profile, len(vars), max)
	}
	return nil
}

// SetOverride sets a per-profile variable limit.
func (m *Manager) SetOverride(profile string, max int) {
	m.policy.Overrides[profile] = max
}

// RemoveOverride removes a per-profile override, reverting to default.
func (m *Manager) RemoveOverride(profile string) {
	delete(m.policy.Overrides, profile)
}
