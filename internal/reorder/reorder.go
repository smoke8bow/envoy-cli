package reorder

import "fmt"

// Store is the minimal interface required by the Reorderer.
type Store interface {
	Get(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// Reorderer provides key-ordering utilities for profiles.
type Reorderer struct {
	store Store
}

// NewReorderer creates a new Reorderer backed by store.
func NewReorderer(s Store) *Reorderer {
	return &Reorderer{store: s}
}

// Apply reorders the keys of profile name according to the provided ordered
// key list. Keys not present in the list are appended at the end in their
// original iteration order. Keys in the list that do not exist in the profile
// are silently ignored.
func (r *Reorderer) Apply(name string, order []string) (map[string]string, error) {
	vars, err := r.store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("reorder: get profile: %w", err)
	}

	seen := make(map[string]bool, len(order))
	result := make(map[string]string, len(vars))

	for _, k := range order {
		if v, ok := vars[k]; ok {
			result[k] = v
			seen[k] = true
		}
	}

	for k, v := range vars {
		if !seen[k] {
			result[k] = v
		}
	}

	if err := r.store.Save(name, result); err != nil {
		return nil, fmt.Errorf("reorder: save profile: %w", err)
	}

	return result, nil
}

// Preview returns what the profile would look like after reordering without
// persisting any changes.
func (r *Reorderer) Preview(name string, order []string) ([]string, error) {
	vars, err := r.store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("reorder: get profile: %w", err)
	}

	seen := make(map[string]bool, len(order))
	var keys []string

	for _, k := range order {
		if _, ok := vars[k]; ok {
			keys = append(keys, k)
			seen[k] = true
		}
	}

	for k := range vars {
		if !seen[k] {
			keys = append(keys, k)
		}
	}

	return keys, nil
}
