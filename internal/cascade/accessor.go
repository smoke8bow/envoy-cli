package cascade

// StoreAccessor adapts any store that exposes Get(name) (map[string]string, error)
// so it satisfies the Accessor interface without importing the store package
// directly (avoids circular dependencies).
type StoreAccessor struct {
	getFn func(name string) (map[string]string, error)
}

// NewStoreAccessor wraps a bare function as an Accessor.
func NewStoreAccessor(fn func(name string) (map[string]string, error)) *StoreAccessor {
	return &StoreAccessor{getFn: fn}
}

// Get implements Accessor.
func (a *StoreAccessor) Get(name string) (map[string]string, error) {
	return a.getFn(name)
}
