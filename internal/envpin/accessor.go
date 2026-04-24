package envpin

import "fmt"

// StoreReader is satisfied by any store that can return profile vars.
type StoreReader interface {
	Get(profile string) (map[string]string, error)
}

// StoreWriter is satisfied by any store that can persist profile vars.
type StoreWriter interface {
	Set(profile string, vars map[string]string) error
}

// StoreAccessor combines read and write access to a profile store.
type StoreAccessor interface {
	StoreReader
	StoreWriter
}

// GuardWrite prevents writing pinned keys in a profile.
// It loads the current vars, strips pinned keys from updates, and saves.
func GuardWrite(m *Manager, store StoreAccessor, profile string, updates map[string]string) error {
	existing, err := store.Get(profile)
	if err != nil {
		return fmt.Errorf("envpin guard: get %q: %w", profile, err)
	}

	// Build merged map: start from existing, apply only non-pinned updates.
	result := make(map[string]string, len(existing))
	for k, v := range existing {
		result[k] = v
	}
	for k, v := range m.FilterWritable(profile, updates) {
		result[k] = v
	}

	if err := store.Set(profile, result); err != nil {
		return fmt.Errorf("envpin guard: set %q: %w", profile, err)
	}
	return nil
}

// GuardDelete prevents deleting pinned keys from a profile.
func GuardDelete(m *Manager, store StoreAccessor, profile, key string) error {
	if m.IsPinned(profile, key) {
		return fmt.Errorf("envpin: cannot delete pinned key %q in profile %q", key, profile)
	}
	existing, err := store.Get(profile)
	if err != nil {
		return fmt.Errorf("envpin guard: get %q: %w", profile, err)
	}
	delete(existing, key)
	if err := store.Set(profile, existing); err != nil {
		return fmt.Errorf("envpin guard: set %q: %w", profile, err)
	}
	return nil
}
