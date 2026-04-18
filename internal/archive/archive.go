package archive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents an archived snapshot of a profile.
type Entry struct {
	Profile string            `json:"profile"`
	Vars    map[string]string `json:"vars"`
	ArchivedAt time.Time     `json:"archived_at"`
}

// Manager handles archiving and restoring profiles.
type Manager struct {
	dir string
}

func NewManager(dir string) *Manager {
	return &Manager{dir: dir}
}

func (m *Manager) Archive(profile string, vars map[string]string) error {
	if err := os.MkdirAll(m.dir, 0700); err != nil {
		return fmt.Errorf("archive: mkdir: %w", err)
	}
	e := Entry{
		Profile:    profile,
		Vars:       vars,
		ArchivedAt: time.Now().UTC(),
	}
	filename := fmt.Sprintf("%s_%d.json", profile, e.ArchivedAt.UnixNano())
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: marshal: %w", err)
	}
	return os.WriteFile(filepath.Join(m.dir, filename), data, 0600)
}

func (m *Manager) List(profile string) ([]Entry, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("archive: readdir: %w", err)
	}
	var results []Entry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(m.dir, e.Name()))
		if err != nil {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		if profile == "" || entry.Profile == profile {
			results = append(results, entry)
		}
	}
	return results, nil
}

func (m *Manager) Latest(profile string) (*Entry, error) {
	all, err := m.List(profile)
	if err != nil {
		return nil, err
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("archive: no entries for profile %q", profile)
	}
	latest := all[0]
	for _, e := range all[1:] {
		if e.ArchivedAt.After(latest.ArchivedAt) {
			latest = e
		}
	}
	return &latest, nil
}
