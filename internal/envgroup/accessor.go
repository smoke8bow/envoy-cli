package envgroup

import "fmt"

// ProfileGetter is satisfied by store.Store and similar types.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// GroupProfile loads a named profile from store and groups its variables.
func GroupProfile(store ProfileGetter, profile string, opts Options) (Result, error) {
	vars, err := store.Get(profile)
	if err != nil {
		return Result{}, fmt.Errorf("envgroup: load profile %q: %w", profile, err)
	}
	return GroupBy(vars, opts), nil
}
