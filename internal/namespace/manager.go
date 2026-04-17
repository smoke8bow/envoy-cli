package namespace

import (
	"errors"
	"sort"
)

// FSManager manages namespaces backed by a FileStore.
type FSManager struct {
	fs *FileStore
}

func NewFSManager(fs *FileStore) *FSManager {
	return &FSManager{fs: fs}
}

func (m *FSManager) List() ([]Namespace, error) {
	return m.fs.Load()
}

func (m *FSManager) Create(name string) error {
	if name == "" {
		return errors.New("namespace name must not be empty")
	}
	ns, err := m.fs.Load()
	if err != nil {
		return err
	}
	for _, n := range ns {
		if n.Name == name {
			return errors.New("namespace already exists: " + name)
		}
	}
	ns = append(ns, Namespace{Name: name, Profiles: []string{}})
	return m.fs.Save(ns)
}

func (m *FSManager) Delete(name string) error {
	ns, err := m.fs.Load()
	if err != nil {
		return err
	}
	for i, n := range ns {
		if n.Name == name {
			ns = append(ns[:i], ns[i+1:]...)
			return m.fs.Save(ns)
		}
	}
	return errors.New("namespace not found: " + name)
}

func (m *FSManager) Assign(nsName, profile string) error {
	ns, err := m.fs.Load()
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
			return m.fs.Save(ns)
		}
	}
	return errors.New("namespace not found: " + nsName)
}

func (m *FSManager) Unassign(nsName, profile string) error {
	ns, err := m.fs.Load()
	if err != nil {
		return err
	}
	for i, n := range ns {
		if n.Name == nsName {
			for j, p := range n.Profiles {
				if p == profile {
					ns[i].Profiles = append(n.Profiles[:j], n.Profiles[j+1:]...)
					return m.fs.Save(ns)
				}
			}
			return errors.New("profile not in namespace")
		}
	}
	return errors.New("namespace not found: " + nsName)
}
