package envrotate

import "fmt"

// StoreAccessor wraps a store that exposes Get/Set via a common interface
// used by other packages in envoy-cli (e.g. internal/store).
type StoreAccessor struct {
	getter func(string) (map[string]string, error)
	setter func(string, map[string]string) error
}

// NewStoreAccessor creates a StoreAccessor from plain function values,
// allowing callers to adapt any store without implementing the full interface.
func NewStoreAccessor(
	getter func(string) (map[string]string, error),
	setter func(string, map[string]string) error,
) *StoreAccessor {
	return &StoreAccessor{getter: getter, setter: setter}
}

func (a *StoreAccessor) Get(name string) (map[string]string, error) {
	return a.getter(name)
}

func (a *StoreAccessor) Set(name string, vars map[string]string) error {
	return a.setter(name, vars)
}

// RotateProfile is a convenience function that constructs a Manager and
// performs a single rotation in one call.
func RotateProfile(
	store Store,
	profile string,
	rotationMap map[string]string,
	opts Options,
) (Result, error) {
	if store == nil {
		return Result{}, fmt.Errorf("envrotate: store must not be nil")
	}
	return NewManager(store, opts).Rotate(profile, rotationMap)
}
