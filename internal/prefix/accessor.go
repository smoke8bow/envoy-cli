package prefix

// ProfileAccessor abstracts the store dependency for profile-level operations.
type ProfileAccessor interface {
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
}

// ApplyToProfile loads a profile, applies the prefix to all keys, and saves it back.
func ApplyToProfile(store ProfileAccessor, profile, p string) error {
	vars, err := store.Get(profile)
	if err != nil {
		return err
	}
	m := NewManager()
	prefixed, err := m.Apply(vars, p)
	if err != nil {
		return err
	}
	return store.Set(profile, prefixed)
}

// StripFromProfile loads a profile, strips the prefix from all keys, and saves it back.
func StripFromProfile(store ProfileAccessor, profile, p string, onlyPrefixed bool) error {
	vars, err := store.Get(profile)
	if err != nil {
		return err
	}
	m := NewManager()
	stripped, err := m.Strip(vars, p, onlyPrefixed)
	if err != nil {
		return err
	}
	return store.Set(profile, stripped)
}
