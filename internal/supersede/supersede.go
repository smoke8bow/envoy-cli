// Package supersede provides functionality to override specific keys in a
// profile with values from another profile, without performing a full merge.
package supersede

import "fmt"

// Store is the minimal interface required to read and write profiles.
type Store interface {
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
}

// Manager applies targeted key overrides from a source profile onto a
// destination profile.
type Manager struct {
	store Store
}

// NewManager creates a new Manager backed by the given store.
func NewManager(store Store) *Manager {
	return &Manager{store: store}
}

// Apply copies the specified keys from src into dst. If keys is empty, all
// keys from src are copied. Keys that do not exist in src are skipped.
// Returns an error if either profile cannot be read or the destination cannot
// be saved.
func (m *Manager) Apply(dst, src string, keys []string) ([]string, error) {
	if dst == "" {
		return nil, fmt.Errorf("supersede: destination profile name must not be empty")
	}
	if src == "" {
		return nil, fmt.Errorf("supersede: source profile name must not be empty")
	}

	dstVars, err := m.store.Get(dst)
	if err != nil {
		return nil, fmt.Errorf("supersede: load destination %q: %w", dst, err)
	}

	srcVars, err := m.store.Get(src)
	if err != nil {
		return nil, fmt.Errorf("supersede: load source %q: %w", src, err)
	}

	targets := keys
	if len(targets) == 0 {
		for k := range srcVars {
			targets = append(targets, k)
		}
	}

	updated := make(map[string]string, len(dstVars))
	for k, v := range dstVars {
		updated[k] = v
	}

	var applied []string
	for _, k := range targets {
		v, ok := srcVars[k]
		if !ok {
			continue
		}
		updated[k] = v
		applied = append(applied, k)
	}

	if err := m.store.Set(dst, updated); err != nil {
		return nil, fmt.Errorf("supersede: save destination %q: %w", dst, err)
	}

	return applied, nil
}
