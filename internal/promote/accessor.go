package promote

// StoreAccessor wraps a concrete store type that satisfies the Store interface,
// allowing callers to build a Promoter without importing the store package
// directly.
type StoreAccessor struct {
	get  func(string) (map[string]string, error)
	save func(string, map[string]string) error
	list func() ([]string, error)
}

// NewStoreAccessor constructs a StoreAccessor from plain function values.
func NewStoreAccessor(
	get func(string) (map[string]string, error),
	save func(string, map[string]string) error,
	list func() ([]string, error),
) *StoreAccessor {
	return &StoreAccessor{get: get, save: save, list: list}
}

func (a *StoreAccessor) Get(name string) (map[string]string, error) {
	return a.get(name)
}

func (a *StoreAccessor) Save(name string, vars map[string]string) error {
	return a.save(name, vars)
}

func (a *StoreAccessor) List() ([]string, error) {
	return a.list()
}
