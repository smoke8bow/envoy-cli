package rollback

import (
	"errors"
	"fmt"
)

// Store is the minimal interface required by the rollback manager.
type Store interface {
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
}

// SnapshotManager is the minimal interface for taking/restoring snapshots.
type SnapshotManager interface {
	Take(profile string) (string, error)
	Restore(snapshotName string) error
}

// Manager provides rollback support for profiles.
type Manager struct {
	store    Store
	snapshot SnapshotManager
}

// NewManager creates a new rollback Manager.
func NewManager(store Store, snapshot SnapshotManager) *Manager {
	return &Manager{store: store, snapshot: snapshot}
}

// Checkpoint takes a snapshot of the profile before a destructive operation.
// Returns the snapshot name so it can be passed to Rollback if needed.
func (m *Manager) Checkpoint(profile string) (string, error) {
	if _, err := m.store.Get(profile); err != nil {
		return "", fmt.Errorf("rollback: profile %q not found: %w", profile, err)
	}
	name, err := m.snapshot.Take(profile)
	if err != nil {
		return "", fmt.Errorf("rollback: checkpoint failed: %w", err)
	}
	return name, nil
}

// Rollback restores a profile to the state captured in the given snapshot.
func (m *Manager) Rollback(snapshotName string) error {
	if snapshotName == "" {
		return errors.New("rollback: snapshot name must not be empty")
	}
	if err := m.snapshot.Restore(snapshotName); err != nil {
		return fmt.Errorf("rollback: restore failed: %w", err)
	}
	return nil
}
