package scope

import (
	"errors"
	"fmt"
	"strings"
)

// Scope represents a named context that can restrict or namespace profile operations.
type Scope struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

// Manager manages a set of named scopes.
type Manager struct {
	scopes map[string]*Scope
}

func NewManager() *Manager {
	return &Manager{scopes: make(map[string]*Scope)}
}

func (m *Manager) Create(name string, labels map[string]string) (*Scope, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("scope name must not be empty")
	}
	if _, exists := m.scopes[name]; exists {
		return nil, fmt.Errorf("scope %q already exists", name)
	}
	lbls := make(map[string]string)
	for k, v := range labels {
		lbls[k] = v
	}
	s := &Scope{Name: name, Labels: lbls}
	m.scopes[name] = s
	return s, nil
}

func (m *Manager) Get(name string) (*Scope, error) {
	s, ok := m.scopes[name]
	if !ok {
		return nil, fmt.Errorf("scope %q not found", name)
	}
	return s, nil
}

func (m *Manager) Delete(name string) error {
	if _, ok := m.scopes[name]; !ok {
		return fmt.Errorf("scope %q not found", name)
	}
	delete(m.scopes, name)
	return nil
}

func (m *Manager) List() []*Scope {
	out := make([]*Scope, 0, len(m.scopes))
	for _, s := range m.scopes {
		out = append(out, s)
	}
	return out
}

// Match returns true if the scope's labels contain all of the given selectors.
func (s *Scope) Match(selectors map[string]string) bool {
	for k, v := range selectors {
		if s.Labels[k] != v {
			return false
		}
	}
	return true
}
