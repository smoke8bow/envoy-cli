// Package envrotate provides key rotation for environment profiles.
// It renames a set of keys within a profile according to a rotation map,
// preserving values and optionally removing the old keys.
package envrotate

import (
	"errors"
	"fmt"
)

// Getter loads a profile's env vars by name.
type Getter interface {
	Get(name string) (map[string]string, error)
}

// Setter persists a profile's env vars by name.
type Setter interface {
	Set(name string, vars map[string]string) error
}

// Store combines Getter and Setter.
type Store interface {
	Getter
	Setter
}

// Options controls rotation behaviour.
type Options struct {
	// RemoveOld removes the old key after copying the value to the new key.
	RemoveOld bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{RemoveOld: true}
}

// Result describes what changed during a rotation.
type Result struct {
	Rotated []Rename
	Skipped []string // old keys that were not present in the profile
}

// Rename records a single key rename.
type Rename struct {
	OldKey string
	NewKey string
}

// Manager performs key rotations against a store.
type Manager struct {
	store Store
	opts  Options
}

// NewManager creates a Manager with the given store and options.
func NewManager(store Store, opts Options) *Manager {
	return &Manager{store: store, opts: opts}
}

// Rotate applies the rotationMap (oldKey -> newKey) to the named profile.
// It returns a Result describing which keys were rotated and which were absent.
func (m *Manager) Rotate(profile string, rotationMap map[string]string) (Result, error) {
	if profile == "" {
		return Result{}, errors.New("envrotate: profile name must not be empty")
	}
	if len(rotationMap) == 0 {
		return Result{}, errors.New("envrotate: rotation map must not be empty")
	}

	vars, err := m.store.Get(profile)
	if err != nil {
		return Result{}, fmt.Errorf("envrotate: load profile %q: %w", profile, err)
	}

	updated := make(map[string]string, len(vars))
	for k, v := range vars {
		updated[k] = v
	}

	var result Result
	for oldKey, newKey := range rotationMap {
		val, ok := updated[oldKey]
		if !ok {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		updated[newKey] = val
		if m.opts.RemoveOld {
			delete(updated, oldKey)
		}
		result.Rotated = append(result.Rotated, Rename{OldKey: oldKey, NewKey: newKey})
	}

	if err := m.store.Set(profile, updated); err != nil {
		return Result{}, fmt.Errorf("envrotate: save profile %q: %w", profile, err)
	}
	return result, nil
}
