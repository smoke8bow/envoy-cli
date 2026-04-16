package pin

import (
	"errors"
	"fmt"

	"envoy-cli/internal/store"
)

const metaKeyPinned = "pinned"

// Manager handles pinning profiles to prevent accidental modification.
type Manager struct {
	store *store.Store
}

func NewManager(s *store.Store) *Manager {
	return &Manager{store: s}
}

// Pin marks a profile as pinned.
func (m *Manager) Pin(name string) error {
	if _, err := m.store.Get(name); err != nil {
		return fmt.Errorf("profile %q not found", name)
	}
	return m.store.SetMeta(name, metaKeyPinned, "true")
}

// Unpin removes the pinned mark from a profile.
func (m *Manager) Unpin(name string) error {
	if _, err := m.store.Get(name); err != nil {
		return fmt.Errorf("profile %q not found", name)
	}
	v, err := m.store.GetMeta(name, metaKeyPinned)
	if err != nil || v != "true" {
		return errors.New("profile is not pinned")
	}
	return m.store.DeleteMeta(name, metaKeyPinned)
}

// IsPinned reports whether a profile is pinned.
func (m *Manager) IsPinned(name string) bool {
	v, err := m.store.GetMeta(name, metaKeyPinned)
	return err == nil && v == "true"
}

// ListPinned returns all pinned profile names.
func (m *Manager) ListPinned() ([]string, error) {
	all, err := m.store.List()
	if err != nil {
		return nil, err
	}
	var pinned []string
	for _, name := range all {
		if m.IsPinned(name) {
			pinned = append(pinned, name)
		}
	}
	return pinned, nil
}
