// Package envpromote provides functionality to promote environment variables
// from one named profile to another, with optional key filtering.
package envpromote

import (
	"errors"
	"fmt"
)

// Getter retrieves a profile's environment variables.
type Getter interface {
	Get(name string) (map[string]string, error)
}

// Setter persists a profile's environment variables.
type Setter interface {
	Set(name string, vars map[string]string) error
}

// Store combines Getter and Setter.
type Store interface {
	Getter
	Setter
}

// Options controls how promotion behaves.
type Options struct {
	// Keys restricts promotion to specific keys. Empty means all keys.
	Keys []string
	// Overwrite controls whether existing keys in dst are replaced.
	Overwrite bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Overwrite: true}
}

// Manager handles profile-to-profile promotion.
type Manager struct {
	store Store
}

// NewManager creates a Manager backed by store.
func NewManager(store Store) *Manager {
	return &Manager{store: store}
}

// Promote copies variables from src profile into dst profile.
func (m *Manager) Promote(src, dst string, opts Options) (map[string]string, error) {
	if src == "" {
		return nil, errors.New("envpromote: src profile name must not be empty")
	}
	if dst == "" {
		return nil, errors.New("envpromote: dst profile name must not be empty")
	}
	if src == dst {
		return nil, errors.New("envpromote: src and dst must be different profiles")
	}

	srcVars, err := m.store.Get(src)
	if err != nil {
		return nil, fmt.Errorf("envpromote: load src %q: %w", src, err)
	}

	dstVars, err := m.store.Get(dst)
	if err != nil {
		return nil, fmt.Errorf("envpromote: load dst %q: %w", dst, err)
	}

	// Clone dst so we don't mutate the caller's map.
	result := make(map[string]string, len(dstVars))
	for k, v := range dstVars {
		result[k] = v
	}

	promoted := selectKeys(srcVars, opts.Keys)
	for k, v := range promoted {
		if _, exists := result[k]; exists && !opts.Overwrite {
			continue
		}
		result[k] = v
	}

	if err := m.store.Set(dst, result); err != nil {
		return nil, fmt.Errorf("envpromote: save dst %q: %w", dst, err)
	}
	return result, nil
}

// selectKeys returns a filtered copy of vars. If keys is empty, all vars are returned.
func selectKeys(vars map[string]string, keys []string) map[string]string {
	if len(keys) == 0 {
		out := make(map[string]string, len(vars))
		for k, v := range vars {
			out[k] = v
		}
		return out
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := vars[k]; ok {
			out[k] = v
		}
	}
	return out
}
