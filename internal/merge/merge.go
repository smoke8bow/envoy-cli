package merge

import (
	"fmt"
)

// Strategy defines how conflicts are resolved when merging profiles.
type Strategy int

const (
	// StrategyOurs keeps the destination value on conflict.
	StrategyOurs Strategy = iota
	// StrategyTheirs overwrites with the source value on conflict.
	StrategyTheirs
	// StrategyError returns an error on any conflict.
	StrategyError
)

// Result holds the merged environment map and metadata about the operation.
type Result struct {
	Vars      map[string]string
	Added     []string
	Overwrite []string
	Skipped   []string
}

// Merge combines src into dst according to the given strategy.
// dst is never mutated; a new map is returned inside Result.
func Merge(dst, src map[string]string, strategy Strategy) (*Result, error) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	result := &Result{Vars: out}

	for k, v := range src {
		if _, exists := out[k]; exists {
			switch strategy {
			case StrategyOurs:
				result.Skipped = append(result.Skipped, k)
				continue
			case StrategyTheirs:
				out[k] = v
				result.Overwrite = append(result.Overwrite, k)
			case StrategyError:
				return nil, fmt.Errorf("merge conflict on key %q", k)
			}
		} else {
			out[k] = v
			result.Added = append(result.Added, k)
		}
	}

	return result, nil
}
