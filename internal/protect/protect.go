// Package protect provides protection rules that prevent accidental
// modification or deletion of critical environment variable keys.
package protect

import (
	"errors"
	"fmt"
	"sync"
)

// ErrKeyProtected is returned when an operation targets a protected key.
var ErrKeyProtected = errors.New("key is protected")

// ErrKeyNotProtected is returned when trying to unprotect a key that isn't protected.
var ErrKeyNotProtected = errors.New("key is not protected")

// Manager manages protected keys per profile.
type Manager struct {
	mu      sync.RWMutex
	records map[string]map[string]struct{} // profile -> set of protected keys
}

// NewManager returns a new Manager.
func NewManager() *Manager {
	return &Manager{
		records: make(map[string]map[string]struct{}),
	}
}

// Protect marks a key as protected within a profile.
func (m *Manager) Protect(profile, key string) error {
	if profile == "" {
		return errors.New("profile name must not be empty")
	}
	if key == "" {
		return errors.New("key must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.records[profile]; !ok {
		m.records[profile] = make(map[string]struct{})
	}
	m.records[profile][key] = struct{}{}
	return nil
}

// Unprotect removes protection from a key within a profile.
func (m *Manager) Unprotect(profile, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys, ok := m.records[profile]
	if !ok {
		return fmt.Errorf("%w: %s", ErrKeyNotProtected, key)
	}
	if _, exists := keys[key]; !exists {
		return fmt.Errorf("%w: %s", ErrKeyNotProtected, key)
	}
	delete(keys, key)
	if len(keys) == 0 {
		delete(m.records, profile)
	}
	return nil
}

// IsProtected reports whether a key is protected within a profile.
func (m *Manager) IsProtected(profile, key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys, ok := m.records[profile]
	if !ok {
		return false
	}
	_, exists := keys[key]
	return exists
}

// List returns all protected keys for a profile.
func (m *Manager) List(profile string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys, ok := m.records[profile]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	return out
}

// Guard returns an error if any of the given keys are protected in the profile.
func (m *Manager) Guard(profile string, keys []string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	protected, ok := m.records[profile]
	if !ok {
		return nil
	}
	for _, k := range keys {
		if _, exists := protected[k]; exists {
			return fmt.Errorf("%w: %s", ErrKeyProtected, k)
		}
	}
	return nil
}
