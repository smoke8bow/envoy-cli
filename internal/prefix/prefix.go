package prefix

import (
	"fmt"
	"strings"
)

// Manager applies or strips a key prefix across a profile's env vars.
type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// Apply returns a new map with all keys prefixed by p.
func (m *Manager) Apply(vars map[string]string, p string) (map[string]string, error) {
	if p == "" {
		return nil, fmt.Errorf("prefix must not be empty")
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[p+k] = v
	}
	return out, nil
}

// Strip returns a new map with the leading prefix p removed from all keys.
// Keys that do not carry the prefix are omitted when onlyPrefixed is true,
// or kept unchanged when false.
func (m *Manager) Strip(vars map[string]string, p string, onlyPrefixed bool) (map[string]string, error) {
	if p == "" {
		return nil, fmt.Errorf("prefix must not be empty")
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if strings.HasPrefix(k, p) {
			out[strings.TrimPrefix(k, p)] = v
		} else if !onlyPrefixed {
			out[k] = v
		}
	}
	return out, nil
}

// Filter returns only the entries whose keys start with p.
func (m *Manager) Filter(vars map[string]string, p string) map[string]string {
	out := make(map[string]string)
	for k, v := range vars {
		if strings.HasPrefix(k, p) {
			out[k] = v
		}
	}
	return out
}

// Replace swaps an existing prefix oldP with newP on all matching keys.
// Keys that do not start with oldP are left unchanged.
func (m *Manager) Replace(vars map[string]string, oldP, newP string) (map[string]string, error) {
	if oldP == "" {
		return nil, fmt.Errorf("old prefix must not be empty")
	}
	if newP == "" {
		return nil, fmt.Errorf("new prefix must not be empty")
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if strings.HasPrefix(k, oldP) {
			out[newP+strings.TrimPrefix(k, oldP)] = v
		} else {
			out[k] = v
		}
	}
	return out, nil
}
