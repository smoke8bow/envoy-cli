package envpromote

import "fmt"

// StoreAccessor adapts a generic key/value store to the Store interface
// expected by Manager. It assumes the underlying store holds
// map[string]string values keyed by profile name.
type StoreAccessor struct {
	getter func(name string) (map[string]string, error)
	setter func(name string, vars map[string]string) error
}

// NewStoreAccessor creates a StoreAccessor from explicit getter/setter funcs.
func NewStoreAccessor(
	getter func(string) (map[string]string, error),
	setter func(string, map[string]string) error,
) *StoreAccessor {
	return &StoreAccessor{getter: getter, setter: setter}
}

// Get implements Store.
func (a *StoreAccessor) Get(name string) (map[string]string, error) {
	v, err := a.getter(name)
	if err != nil {
		return nil, fmt.Errorf("accessor: get %q: %w", name, err)
	}
	return v, nil
}

// Set implements Store.
func (a *StoreAccessor) Set(name string, vars map[string]string) error {
	if err := a.setter(name, vars); err != nil {
		return fmt.Errorf("accessor: set %q: %w", name, err)
	}
	return nil
}
