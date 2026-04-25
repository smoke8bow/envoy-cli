package envflag

import "fmt"

// ProfileStore is satisfied by the project's standard store.
type ProfileStore interface {
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
}

// storeAdapter wraps a ProfileStore to satisfy the Store interface used by
// Manager. It silently returns an empty map when a profile has not been
// created yet, matching the behaviour of other accessor helpers in the
// project.
type storeAdapter struct {
	inner ProfileStore
}

func (a *storeAdapter) Get(name string) (map[string]string, error) {
	vars, err := a.inner.Get(name)
	if err != nil {
		return map[string]string{}, nil
	}
	return vars, nil
}

func (a *storeAdapter) Set(name string, vars map[string]string) error {
	return a.inner.Set(name, vars)
}

// NewStoreAccessor returns a Manager that reads from and writes to the
// provided ProfileStore, using profile as the backing key.
func NewStoreAccessor(store ProfileStore, profile string) (*Manager, error) {
	if profile == "" {
		return nil, fmt.Errorf("envflag: accessor profile must not be empty")
	}
	return NewManager(&storeAdapter{inner: store}, profile)
}
