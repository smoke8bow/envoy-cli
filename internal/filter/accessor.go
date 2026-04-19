package filter

import "fmt"

// StoreAccessor abstracts reading a profile's vars from a store.
type StoreAccessor interface {
	Get(name string) (map[string]string, error)
}

// FilterProfile loads a profile by name and applies the given Option.
// It returns the filtered Result or an error if the profile cannot be loaded.
func FilterProfile(store StoreAccessor, profile string, opt Option) (Result, error) {
	vars, err := store.Get(profile)
	if err != nil {
		return Result{}, fmt.Errorf("filter: load profile %q: %w", profile, err)
	}
	return Filter(vars, opt), nil
}
