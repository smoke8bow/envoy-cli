package chain

import "fmt"

// ProfileGetter is the minimal interface needed to fetch a profile's vars.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// FromProfiles builds a Composer by loading each named profile from the store.
// Profiles are applied in the order provided.
func FromProfiles(store ProfileGetter, names []string) (*Composer, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("chain: at least one profile name is required")
	}
	entries := make([]Entry, 0, len(names))
	for _, n := range names {
		vars, err := store.Get(n)
		if err != nil {
			return nil, fmt.Errorf("chain: loading profile %q: %w", n, err)
		}
		entries = append(entries, Entry{Name: n, Vars: vars})
	}
	return NewComposer(entries)
}
