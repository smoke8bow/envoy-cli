// Package draft provides a staging area for building up environment variable
// sets before committing them to a named profile.
package draft

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNoDraft is returned when no draft exists for the given name.
var ErrNoDraft = errors.New("no draft found")

// Draft holds a set of pending key/value changes.
type Draft struct {
	Vars map[string]string
}

// Manager manages in-memory drafts keyed by a draft name.
type Manager struct {
	mu     sync.Mutex
	drafts map[string]*Draft
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	return &Manager{drafts: make(map[string]*Draft)}
}

// Create initialises an empty draft with the given name.
// Returns an error if a draft with that name already exists.
func (m *Manager) Create(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.drafts[name]; ok {
		return fmt.Errorf("draft %q already exists", name)
	}
	m.drafts[name] = &Draft{Vars: make(map[string]string)}
	return nil
}

// Set adds or updates a key in the named draft.
func (m *Manager) Set(name, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.drafts[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoDraft, name)
	}
	d.Vars[key] = value
	return nil
}

// Delete removes a key from the named draft.
func (m *Manager) Delete(name, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.drafts[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoDraft, name)
	}
	delete(d.Vars, key)
	return nil
}

// Get returns a copy of the draft vars for the given name.
func (m *Manager) Get(name string) (map[string]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.drafts[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoDraft, name)
	}
	out := make(map[string]string, len(d.Vars))
	for k, v := range d.Vars {
		out[k] = v
	}
	return out, nil
}

// Discard removes the draft entirely.
func (m *Manager) Discard(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.drafts[name]; !ok {
		return fmt.Errorf("%w: %s", ErrNoDraft, name)
	}
	delete(m.drafts, name)
	return nil
}

// List returns all current draft names.
func (m *Manager) List() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	names := make([]string, 0, len(m.drafts))
	for n := range m.drafts {
		names = append(names, n)
	}
	return names
}
