package envscope

import "fmt"

// StoreAccessor adapts a map-based store to the Getter interface.
type StoreAccessor struct {
	data map[string]map[string]string
}

// NewStoreAccessor wraps an in-memory map as a Getter.
func NewStoreAccessor(data map[string]map[string]string) *StoreAccessor {
	return &StoreAccessor{data: data}
}

// Get returns the vars for the named profile or an error if not found.
func (s *StoreAccessor) Get(name string) (map[string]string, error) {
	vars, ok := s.data[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return copy, nil
}
