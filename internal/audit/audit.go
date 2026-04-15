package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventCreate EventType = "create"
	EventUpdate EventType = "update"
	EventDelete EventType = "delete"
	EventSwitch EventType = "switch"
)

// Entry records a single audit event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Profile   string    `json:"profile"`
	Detail    string    `json:"detail,omitempty"`
}

// Log holds a list of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
	path    string
}

// Load reads the audit log from disk, returning an empty log if the file does not exist.
func Load(dir string) (*Log, error) {
	p := filepath.Join(dir, "audit.json")
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &Log{path: p}, nil
	}
	if err != nil {
		return nil, err
	}
	var l Log
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, err
	}
	l.path = p
	return &l, nil
}

// Record appends a new entry to the log and persists it.
func (l *Log) Record(event EventType, profile, detail string) error {
	l.Entries = append(l.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Profile:   profile,
		Detail:    detail,
	})
	return l.save()
}

// Recent returns up to n most recent entries.
func (l *Log) Recent(n int) []Entry {
	if n >= len(l.Entries) {
		return l.Entries
	}
	return l.Entries[len(l.Entries)-n:]
}

func (l *Log) save() error {
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(l.path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o600)
}
