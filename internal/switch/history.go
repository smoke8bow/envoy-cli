package switch_

import (
	"encoding/json"
	"os"
	"time"
)

// HistoryEntry records a single profile switch event.
type HistoryEntry struct {
	Profile   string    `json:"profile"`
	SwitchedAt time.Time `json:"switched_at"`
}

// History manages a persistent log of profile switches.
type History struct {
	path    string
	entries []HistoryEntry
}

// LoadHistory reads switch history from path, creating an empty history if the
// file does not yet exist.
func LoadHistory(path string) (*History, error) {
	h := &History{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &h.entries); err != nil {
		return nil, err
	}
	return h, nil
}

// Record appends a new entry for the given profile and persists the history.
func (h *History) Record(profile string) error {
	h.entries = append(h.entries, HistoryEntry{
		Profile:    profile,
		SwitchedAt: time.Now().UTC(),
	})
	return h.save()
}

// Last returns the most recent history entry, or nil if history is empty.
func (h *History) Last() *HistoryEntry {
	if len(h.entries) == 0 {
		return nil
	}
	e := h.entries[len(h.entries)-1]
	return &e
}

// Entries returns all recorded history entries in chronological order.
func (h *History) Entries() []HistoryEntry {
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o600)
}
