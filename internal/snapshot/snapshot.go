package snapshot

import (
	"fmt"
	"time"

	"envoy-cli/internal/store"
)

// Entry represents a point-in-time capture of a profile's environment variables.
type Entry struct {
	ID        string            `json:"id"`
	Profile   string            `json:"profile"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
	Note      string            `json:"note,omitempty"`
}

// Manager handles snapshot creation and retrieval.
type Manager struct {
	store snapshotStore
}

type snapshotStore interface {
	Get(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
}

// NewManager creates a new snapshot Manager backed by the given store.
func NewManager(s snapshotStore) *Manager {
	return &Manager{store: s}
}

// Take captures the current state of a profile and saves it as a new snapshot profile.
// The snapshot is saved under a generated name: "<profile>__snap_<timestamp>".
func (m *Manager) Take(profile, note string) (*Entry, error) {
	vars, err := m.store.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("snapshot: profile %q not found: %w", profile, err)
	}

	now := time.Now().UTC()
	id := fmt.Sprintf("%s__snap_%d", profile, now.UnixMilli())

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	if err := m.store.Save(id, copy); err != nil {
		return nil, fmt.Errorf("snapshot: failed to save: %w", err)
	}

	return &Entry{
		ID:        id,
		Profile:   profile,
		Vars:      copy,
		CreatedAt: now,
		Note:      note,
	}, nil
}

// Restore loads a snapshot entry and writes its vars back to the original profile.
func (m *Manager) Restore(snapshotID, targetProfile string) error {
	vars, err := m.store.Get(snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot: snapshot %q not found: %w", snapshotID, err)
	}
	return m.store.Save(targetProfile, vars)
}

// storeAdapter adapts *store.Store to snapshotStore.
type storeAdapter struct{ s *store.Store }

func (a *storeAdapter) Get(name string) (map[string]string, error) { return a.s.Get(name) }
func (a *storeAdapter) Save(name string, vars map[string]string) error {
	return a.s.Save(name, vars)
}

// NewManagerFromStore creates a Manager from the real store.
func NewManagerFromStore(s *store.Store) *Manager {
	return NewManager(&storeAdapter{s: s})
}
