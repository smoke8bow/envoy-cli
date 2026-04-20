// Package cascade provides ordered multi-profile variable resolution.
// Variables are resolved by walking a chain of profiles from lowest to
// highest priority; later profiles override earlier ones.
package cascade

import "fmt"

// Accessor is the interface required to fetch a profile's variables.
type Accessor interface {
	Get(name string) (map[string]string, error)
}

// Manager resolves environment variables across an ordered list of profiles.
type Manager struct {
	store Accessor
}

// NewManager returns a Manager backed by the given Accessor.
func NewManager(store Accessor) *Manager {
	return &Manager{store: store}
}

// Result holds the merged variables and a per-key source map.
type Result struct {
	Vars   map[string]string
	Source map[string]string // key -> profile name that provided it
}

// Resolve merges profiles in order; later entries take precedence.
// profiles must contain at least one entry.
func (m *Manager) Resolve(profiles []string) (*Result, error) {
	if len(profiles) == 0 {
		return nil, fmt.Errorf("cascade: at least one profile required")
	}

	result := &Result{
		Vars:   make(map[string]string),
		Source: make(map[string]string),
	}

	for _, name := range profiles {
		vars, err := m.store.Get(name)
		if err != nil {
			return nil, fmt.Errorf("cascade: loading profile %q: %w", name, err)
		}
		for k, v := range vars {
			result.Vars[k] = v
			result.Source[k] = name
		}
	}

	return result, nil
}

// Keys returns all resolved keys in deterministic (sorted) order.
func (r *Result) Keys() []string {
	keys := make([]string, 0, len(r.Vars))
	for k := range r.Vars {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}
