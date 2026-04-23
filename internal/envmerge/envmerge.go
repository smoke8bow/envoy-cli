// Package envmerge provides utilities for merging multiple profiles
// with configurable conflict resolution and source tracking.
package envmerge

import (
	"errors"
	"fmt"
	"sort"
)

// Strategy controls how key conflicts are resolved during a merge.
type Strategy string

const (
	StrategyFirst  Strategy = "first"  // keep value from the first source that defines the key
	StrategyLast   Strategy = "last"   // keep value from the last source that defines the key
	StrategyStrict Strategy = "strict" // return an error on any conflict
)

// Supported returns all valid strategy names.
func Supported() []string {
	return []string{string(StrategyFirst), string(StrategyLast), string(StrategyStrict)}
}

// IsSupported reports whether s is a valid Strategy.
func IsSupported(s string) bool {
	for _, v := range Supported() {
		if v == s {
			return true
		}
	}
	return false
}

// Source pairs a named label with its environment variable map.
type Source struct {
	Name string
	Vars map[string]string
}

// Result holds the merged variable map and metadata about each key's origin.
type Result struct {
	Vars   map[string]string
	Origin map[string]string // key → source name that won
}

// Merge combines sources according to the given strategy.
// Sources are processed in order; the strategy determines conflict behaviour.
func Merge(sources []Source, strategy Strategy) (*Result, error) {
	if !IsSupported(string(strategy)) {
		return nil, fmt.Errorf("envmerge: unsupported strategy %q", strategy)
	}

	result := &Result{
		Vars:   make(map[string]string),
		Origin: make(map[string]string),
	}

	for _, src := range sources {
		// Iterate keys in deterministic order.
		keys := make([]string, 0, len(src.Vars))
		for k := range src.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := src.Vars[k]
			if existing, exists := result.Vars[k]; exists {
				switch strategy {
				case StrategyFirst:
					// keep existing — do nothing
					_ = existing
				case StrategyLast:
					result.Vars[k] = v
					result.Origin[k] = src.Name
				case StrategyStrict:
					return nil, errors.New("envmerge: conflict on key " + k +
						" between " + result.Origin[k] + " and " + src.Name)
				}
			} else {
				result.Vars[k] = v
				result.Origin[k] = src.Name
			}
		}
	}

	return result, nil
}
