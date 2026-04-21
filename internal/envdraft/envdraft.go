// Package envdraft provides a temporary staging area for environment variable
// edits before they are committed to a named profile.
package envdraft

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNoDraft is returned when no draft exists for the requested profile.
var ErrNoDraft = errors.New("no draft found")

// ErrDraftExists is returned when a draft already exists for the profile.
var ErrDraftExists = errors.New("draft already exists")

// Draft holds a staged set of environment variable edits.
type Draft struct {
	Profile string
	Vars    map[string]string
}

// Store is the backing interface for committing draft vars.
type Store interface {
	Get(profile string) (map[string]string, error)
	Save(profile string, vars map[string]string) error
}

// Manager manages in-memory drafts keyed by profile name.
type Manager struct {
	mu     sync.Mutex
	drafts map[string]*Draft
	store  Store
}

// NewManager creates a new draft Manager backed by the given Store.
func NewManager(s Store) *Manager {
	return &Manager{
		drafts: make(map[string]*Draft),
		store:  s,
	}
}

// Open creates a new draft for the named profile, seeding it from the store.
// Returns ErrDraftExists if a draft is already open.
func (m *Manager) Open(profile string) (*Draft, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.drafts[profile]; ok {
		return nil, ErrDraftExists
	}

	vars, err := m.store.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("open draft: %w", err)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	d := &Draft{Profile: profile, Vars: copy}
	m.drafts[profile] = d
	return d, nil
}

// Get returns the open draft for a profile or ErrNoDraft.
func (m *Manager) Get(profile string) (*Draft, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	d, ok := m.drafts[profile]
	if !ok {
		return nil, ErrNoDraft
	}
	return d, nil
}

// Set updates a key in the open draft.
func (m *Manager) Set(profile, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	d, ok := m.drafts[profile]
	if !ok {
		return ErrNoDraft
	}
	d.Vars[key] = value
	return nil
}

// Delete removes a key from the open draft.
func (m *Manager) Delete(profile, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	d, ok := m.drafts[profile]
	if !ok {
		return ErrNoDraft
	}
	delete(d.Vars, key)
	return nil
}

// Commit writes the draft vars to the store and closes the draft.
func (m *Manager) Commit(profile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	d, ok := m.drafts[profile]
	if !ok {
		return ErrNoDraft
	}

	if err := m.store.Save(profile, d.Vars); err != nil {
		return fmt.Errorf("commit draft: %w", err)
	}

	delete(m.drafts, profile)
	return nil
}

// Discard closes the draft without saving.
func (m *Manager) Discard(profile string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.drafts[profile]; !ok {
		return ErrNoDraft
	}
	delete(m.drafts, profile)
	return nil
}
