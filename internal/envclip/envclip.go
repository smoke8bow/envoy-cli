// Package envclip provides clipboard-style copy/paste for environment variable subsets.
package envclip

import (
	"errors"
	"fmt"
	"sort"
)

// ErrClipboardEmpty is returned when the clipboard has no content.
var ErrClipboardEmpty = errors.New("clipboard is empty")

// ErrKeyNotFound is returned when a requested key is absent from the clipboard.
var ErrKeyNotFound = errors.New("key not found in clipboard")

// Clipboard holds a named snapshot of environment variable key/value pairs.
type Clipboard struct {
	Source string
	vars   map[string]string
}

// Manager maintains a single in-process clipboard.
type Manager struct {
	clip *Clipboard
}

// NewManager returns a new Manager with an empty clipboard.
func NewManager() *Manager {
	return &Manager{}
}

// Copy stores the given vars under the source profile name.
// Passing an empty keys slice copies all vars.
func (m *Manager) Copy(source string, vars map[string]string, keys []string) error {
	if source == "" {
		return errors.New("source name must not be empty")
	}
	selected := make(map[string]string, len(vars))
	if len(keys) == 0 {
		for k, v := range vars {
			selected[k] = v
		}
	} else {
		for _, k := range keys {
			v, ok := vars[k]
			if !ok {
				return fmt.Errorf("%w: %s", ErrKeyNotFound, k)
			}
			selected[k] = v
		}
	}
	m.clip = &Clipboard{Source: source, vars: selected}
	return nil
}

// Paste merges clipboard contents into dst.
// Existing keys in dst are overwritten.
func (m *Manager) Paste(dst map[string]string) (map[string]string, error) {
	if m.clip == nil {
		return nil, ErrClipboardEmpty
	}
	out := make(map[string]string, len(dst)+len(m.clip.vars))
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range m.clip.vars {
		out[k] = v
	}
	return out, nil
}

// Keys returns the sorted list of keys currently held in the clipboard.
func (m *Manager) Keys() ([]string, error) {
	if m.clip == nil {
		return nil, ErrClipboardEmpty
	}
	keys := make([]string, 0, len(m.clip.vars))
	for k := range m.clip.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

// Clear empties the clipboard.
func (m *Manager) Clear() {
	m.clip = nil
}

// IsEmpty reports whether the clipboard holds any content.
func (m *Manager) IsEmpty() bool {
	return m.clip == nil
}
