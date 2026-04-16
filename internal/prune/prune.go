package prune

import "fmt"

// Store is the interface required by the pruner.
type Store interface {
	List() ([]string, error)
	Get(name string) (map[string]string, error)
	Delete(name string) error
}

// Result holds the outcome of a prune operation.
type Result struct {
	Removed []string
	Skipped []string
}

// Manager handles pruning of profiles.
type Manager struct {
	store Store
}

// NewManager creates a new prune Manager.
func NewManager(s Store) *Manager {
	return &Manager{store: s}
}

// DryRun returns profiles that would be removed (empty var maps).
func (m *Manager) DryRun() ([]string, error) {
	names, err := m.store.List()
	if err != nil {
		return nil, fmt.Errorf("prune: list profiles: %w", err)
	}
	var candidates []string
	for _, name := range names {
		vars, err := m.store.Get(name)
		if err != nil {
			return nil, fmt.Errorf("prune: get profile %q: %w", name, err)
		}
		if len(vars) == 0 {
			candidates = append(candidates, name)
		}
	}
	return candidates, nil
}

// Run removes all profiles with no environment variables.
func (m *Manager) Run() (*Result, error) {
	candidates, err := m.DryRun()
	if err != nil {
		return nil, err
	}
	res := &Result{}
	for _, name := range candidates {
		if err := m.store.Delete(name); err != nil {
			res.Skipped = append(res.Skipped, name)
			continue
		}
		res.Removed = append(res.Removed, name)
	}
	return res, nil
}
