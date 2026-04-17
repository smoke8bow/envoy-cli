package expire

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var ErrExpired = errors.New("profile has expired")

type Entry struct {
	Profile   string    `json:"profile"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Manager struct {
	path string
	entries map[string]Entry
}

func NewManager(dir string) (*Manager, error) {
	m := &Manager{
		path:    filepath.Join(dir, "expiry.json"),
		entries: make(map[string]Entry),
	}
	if err := m.load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) Set(profile string, ttl time.Duration) error {
	m.entries[profile] = Entry{
		Profile:   profile,
		ExpiresAt: time.Now().Add(ttl),
	}
	return m.save()
}

func (m *Manager) Check(profile string) error {
	e, ok := m.entries[profile]
	if !ok {
		return nil
	}
	if time.Now().After(e.ExpiresAt) {
		return ErrExpired
	}
	return nil
}

func (m *Manager) Clear(profile string) error {
	delete(m.entries, profile)
	return m.save()
}

// Purge removes all expired entries and persists the result.
func (m *Manager) Purge() error {
	now := time.Now()
	for profile, e := range m.entries {
		if now.After(e.ExpiresAt) {
			delete(m.entries, profile)
		}
	}
	return m.save()
}

func (m *Manager) load() error {
	data, err := os.ReadFile(m.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.entries)
}

func (m *Manager) save() error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0o644)
}
