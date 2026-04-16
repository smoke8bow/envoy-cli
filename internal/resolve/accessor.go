package resolve

// ProfileAccessor is the minimal interface the Resolver needs to read a
// profile's variables. Other packages may implement this to avoid a direct
// dependency on the store.
type ProfileAccessor interface {
	Get(name string) (map[string]string, error)
}

// ResolveProfile is a convenience helper that loads vars via a ProfileAccessor
// and returns the fully-expanded map.
func ResolveProfile(pa ProfileAccessor, name string, ambient map[string]string) (map[string]string, error) {
	vars, err := pa.Get(name)
	if err != nil {
		return nil, err
	}
	r := NewResolver(ambient)
	return r.Resolve(vars)
}
