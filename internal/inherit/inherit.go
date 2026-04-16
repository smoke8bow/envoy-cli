package inherit

// Inheritor merges a parent profile's vars into a child, without overwriting
// keys already present in the child.
type Inheritor struct {
	store Store
}

// Store is the minimal interface required by Inheritor.
type Store interface {
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
}

// NewInheritor returns an Inheritor backed by the given store.
func NewInheritor(s Store) *Inheritor {
	return &Inheritor{store: s}
}

// Apply copies keys from parent into child for any key not already defined in
// child. Returns the resulting merged map (child is not mutated in the store
// until Commit is called).
func (i *Inheritor) Apply(parent, child string) (map[string]string, error) {
	parentVars, err := i.store.Get(parent)
	if err != nil {
		return nil, fmt.Errorf("inherit: parent %q: %w", parent, err)
	}
	childVars, err := i.store.Get(child)
	if err != nil {
		return nil, fmt.Errorf("inherit: child %q: %w", child, err)
	}

	result := make(map[string]string, len(childVars))
	for k, v := range childVars {
		result[k] = v
	}
	for k, v := range parentVars {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	return result, nil
}

// Commit applies the inheritance and persists the result into the child profile.
func (i *Inheritor) Commit(parent, child string) (map[string]string, error) {
	merged, err := i.Apply(parent, child)
	if err != nil {
		return nil, err
	}
	if err := i.store.Set(child, merged); err != nil {
		return nil, fmt.Errorf("inherit: commit %q: %w", child, err)
	}
	return merged, nil
}
