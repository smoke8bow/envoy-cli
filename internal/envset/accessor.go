package envset

import "fmt"

// ProfileGetter retrieves a profile's variables by name.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// ApplyToProfiles loads two profiles by name, applies op, and returns the result.
func ApplyToProfiles(store ProfileGetter, op Op, profileA, profileB string) (map[string]string, error) {
	if !IsSupported(op) {
		return nil, fmt.Errorf("envset: unsupported operation %q", op)
	}
	a, err := store.Get(profileA)
	if err != nil {
		return nil, fmt.Errorf("envset: loading profile %q: %w", profileA, err)
	}
	b, err := store.Get(profileB)
	if err != nil {
		return nil, fmt.Errorf("envset: loading profile %q: %w", profileB, err)
	}
	return Apply(op, a, b)
}
