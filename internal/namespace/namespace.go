package namespace

import (
	"errors"
	"sort"

	"github.com/envoy-cli/envoy/internal/store"
)

// Namespace groups profiles under a named scope.
type Namespace struct {
	Name     string   `json:"name"`
	Profiles []string `json:"profiles"`
}

// Manager manages namespaces.
type Manager struct {
	store storeIface
}

type storeIface interface {
	Load() (*store.Store, error)
	Save(*store.Store) error
}

const metaKey = "__namespaces__"

func NewManager(s storeIface) *Manager {
	return &Manager{store: s}
}

func (m *Manager) List() ([]Namespace, error) {
	return m.load()
}

func (m *Manager) Create(name string) error {
	if name == "" {
		return errors.New("namespace name must not be empty")
	}
	ns, err := m.load()
	if err != nil {
		return err
	}
	for _, n := range ns {
		if n.Name == name {
			return errors.New("namespace already exists: " + name)
		}
	}
	ns = append(ns, Namespace{Name: name})
	return m.save(ns)
}

func (m *Manager) Delete(name string) error {
	ns, err := m.load()
	if err != nil {
		return err
	}
	idx := -1
	for i, n := range ns {
		if n.Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("namespace not found: " + name)
	}
	ns = append(ns[:idx], ns[idx+1:]...)
	return m.save(ns)
}

func (m *Manager) Assign(nsName, profile string) error {
	ns, err := m.load()
	if err != nil {
		return err
	}
	for i, n := range ns {
		if n.Name == nsName {
			for _, p := range n.Profiles {
				if p == profile {
					return nil
				}
			}
			ns[i].Profiles = append(ns[i].Profiles, profile)
			sort.Strings(ns[i].Profiles)
			return m.save(ns)
		}
	}
	return errors.New("namespace not found: " + nsName)
}

func (m *Manager) load() ([]Namespace, error) {
	return []Namespace{}, nil
}

func (m *Manager) save(_ []Namespace) error {
	return nil
}
