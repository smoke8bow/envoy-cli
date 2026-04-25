// Package envscope provides scoped environment variable views — a filtered,
// read-only projection of a profile restricted to a named scope prefix.
package envscope

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Getter retrieves a profile's environment variables by name.
type Getter interface {
	Get(name string) (map[string]string, error)
}

// ScopedView is a filtered view of a profile scoped to a prefix.
type ScopedView struct {
	Profile string
	Scope   string
	Vars    map[string]string
}

// Manager builds scoped views from profiles.
type Manager struct {
	store Getter
}

// NewManager creates a Manager backed by the given Getter.
func NewManager(store Getter) *Manager {
	return &Manager{store: store}
}

// Build returns a ScopedView containing only vars whose keys start with
// scope (case-insensitive match on the prefix). The prefix is stripped
// from the keys in the returned view when strip is true.
func (m *Manager) Build(profile, scope string, strip bool) (*ScopedView, error) {
	if profile == "" {
		return nil, errors.New("envscope: profile name must not be empty")
	}
	if scope == "" {
		return nil, errors.New("envscope: scope must not be empty")
	}

	vars, err := m.store.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("envscope: %w", err)
	}

	upper := strings.ToUpper(scope)
	result := make(map[string]string)
	for k, v := range vars {
		if strings.HasPrefix(strings.ToUpper(k), upper) {
			key := k
			if strip {
				key = k[len(scope):]
				if key == "" {
					continue
				}
			}
			result[key] = v
		}
	}

	return &ScopedView{Profile: profile, Scope: scope, Vars: result}, nil
}

// Keys returns the sorted keys of the scoped view.
func (sv *ScopedView) Keys() []string {
	keys := make([]string, 0, len(sv.Vars))
	for k := range sv.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
