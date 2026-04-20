package protect

import "fmt"

// StoreReader is the minimal interface needed to read profile vars.
type StoreReader interface {
	Get(profile string) (map[string]string, error)
}

// GuardWrite checks that none of the keys being written to a profile are
// protected. Call this before any set/delete operation on a profile's vars.
func GuardWrite(m *Manager, store StoreReader, profile string, keys []string) error {
	if err := m.Guard(profile, keys); err != nil {
		return fmt.Errorf("protect: write blocked on profile %q: %w", profile, err)
	}
	return nil
}

// GuardDelete checks that none of the keys being deleted from a profile are
// protected.
func GuardDelete(m *Manager, profile string, keys []string) error {
	if err := m.Guard(profile, keys); err != nil {
		return fmt.Errorf("protect: delete blocked on profile %q: %w", profile, err)
	}
	return nil
}
