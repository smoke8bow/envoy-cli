package index

import (
	"fmt"
	"sort"
)

// Entry represents a single index entry mapping a key to the profiles that contain it.
type Entry struct {
	Key      string
	Profiles []string
}

// Index maps environment variable keys to the set of profile names that define them.
type Index map[string]map[string]struct{}

// Builder constructs an Index from a profile store accessor.
type Builder struct {
	loader ProfileLoader
}

// ProfileLoader is the interface required to build an index.
type ProfileLoader interface {
	List() ([]string, error)
	Get(name string) (map[string]string, error)
}

// NewBuilder returns a new Builder backed by the given loader.
func NewBuilder(loader ProfileLoader) *Builder {
	return &Builder{loader: loader}
}

// Build constructs and returns the full index.
func (b *Builder) Build() (Index, error) {
	profiles, err := b.loader.List()
	if err != nil {
		return nil, fmt.Errorf("index: list profiles: %w", err)
	}

	idx := make(Index)
	for _, name := range profiles {
		vars, err := b.loader.Get(name)
		if err != nil {
			return nil, fmt.Errorf("index: get profile %q: %w", name, err)
		}
		for key := range vars {
			if idx[key] == nil {
				idx[key] = make(map[string]struct{})
			}
			idx[key][name] = struct{}{}
		}
	}
	return idx, nil
}

// Lookup returns the sorted list of profile names that contain the given key.
func (idx Index) Lookup(key string) []string {
	set, ok := idx[key]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Entries returns all index entries sorted by key.
func (idx Index) Entries() []Entry {
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{Key: k, Profiles: idx.Lookup(k)})
	}
	return entries
}
