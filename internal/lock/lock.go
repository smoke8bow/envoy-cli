package lock

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrLocked is returned when a profile is already locked.
var ErrLocked = errors.New("profile is already locked")

// ErrNotLocked is returned when attempting to unlock a profile that is not locked.
var ErrNotLocked = errors.New("profile is not locked")

// Manager manages profile locks stored as marker files.
type Manager struct {
	dir string
}

// NewManager creates a new lock Manager using the given directory.
func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

// Lock marks the given profile as locked. Returns ErrLocked if already locked.
func (m *Manager) Lock(profile string) error {
	if m.IsLocked(profile) {
		return ErrLocked
	}
	if err := os.MkdirAll(m.dir, 0o700); err != nil {
		return err
	}
	f, err := os.Create(m.path(profile))
	if err != nil {
		return err
	}
	return f.Close()
}

// Unlock removes the lock for the given profile. Returns ErrNotLocked if not locked.
func (m *Manager) Unlock(profile string) error {
	if !m.IsLocked(profile) {
		return ErrNotLocked
	}
	return os.Remove(m.path(profile))
}

// IsLocked reports whether the given profile is currently locked.
func (m *Manager) IsLocked(profile string) bool {
	_, err := os.Stat(m.path(profile))
	return err == nil
}

// LockedAt returns the time the lock was created, or zero time if not locked.
func (m *Manager) LockedAt(profile string) time.Time {
	info, err := os.Stat(m.path(profile))
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// List returns the names of all currently locked profiles.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (m *Manager) path(profile string) string {
	return filepath.Join(m.dir, profile)
}
