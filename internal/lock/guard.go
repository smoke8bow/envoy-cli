package lock

import "fmt"

// Guard wraps a Manager and provides a checked execution helper.
type Guard struct {
	m *Manager
}

// NewGuard creates a Guard backed by the given Manager.
func NewGuard(m *Manager) *Guard {
	return &Guard{m: m}
}

// Require returns an error if the named profile is locked, preventing mutation.
func (g *Guard) Require(profile string) error {
	if g.m.IsLocked(profile) {
		return fmt.Errorf("profile %q is locked and cannot be modified; unlock it first", profile)
	}
	return nil
}

// WithUnlocked executes fn only if the profile is not locked.
// Returns an error describing the lock if the profile is locked.
func (g *Guard) WithUnlocked(profile string, fn func() error) error {
	if err := g.Require(profile); err != nil {
		return err
	}
	return fn()
}

// Status returns a human-readable status string for the profile.
func (g *Guard) Status(profile string) string {
	if !g.m.IsLocked(profile) {
		return fmt.Sprintf("profile %q is unlocked", profile)
	}
	at := g.m.LockedAt(profile)
	return fmt.Sprintf("profile %q is locked (since %s)", profile, at.Format("2006-01-02 15:04:05"))
}
