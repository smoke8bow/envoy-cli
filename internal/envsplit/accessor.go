package envsplit

import "fmt"

// ProfileGetter retrieves a profile's env vars by name.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// SplitProfile loads a named profile from store and splits it according to rules.
func SplitProfile(store ProfileGetter, profile string, rules []Rule, opts Options) (Result, error) {
	vars, err := store.Get(profile)
	if err != nil {
		return Result{}, fmt.Errorf("envsplit: load profile %q: %w", profile, err)
	}
	return Split(vars, rules, opts)
}
