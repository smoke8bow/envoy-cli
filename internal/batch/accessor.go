package batch

import "fmt"

// StoreAccessor adapts a concrete store to the Store interface expected by Processor.
type StoreAccessor struct {
	get  func(string) (map[string]string, error)
	save func(string, map[string]string) error
}

// NewStoreAccessor creates a StoreAccessor from plain function values.
func NewStoreAccessor(
	get func(string) (map[string]string, error),
	save func(string, map[string]string) error,
) *StoreAccessor {
	if get == nil || save == nil {
		panic("batch: get and save functions must not be nil")
	}
	return &StoreAccessor{get: get, save: save}
}

func (a *StoreAccessor) Get(profile string) (map[string]string, error) {
	v, err := a.get(profile)
	if err != nil {
		return nil, fmt.Errorf("accessor get: %w", err)
	}
	return v, nil
}

func (a *StoreAccessor) Save(profile string, vars map[string]string) error {
	if err := a.save(profile, vars); err != nil {
		return fmt.Errorf("accessor save: %w", err)
	}
	return nil
}
