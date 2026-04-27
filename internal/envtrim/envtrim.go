// Package envtrim removes leading and trailing whitespace from environment
// variable keys and/or values within a profile map.
package envtrim

import (
	"errors"
	"strings"
)

// Options controls which parts of each entry are trimmed.
type Options struct {
	// TrimKeys trims whitespace from keys.
	TrimKeys bool
	// TrimValues trims whitespace from values.
	TrimValues bool
}

// DefaultOptions returns an Options that trims both keys and values.
func DefaultOptions() Options {
	return Options{TrimKeys: true, TrimValues: true}
}

// Result holds the outcome of a Trim call.
type Result struct {
	// Vars is the trimmed environment map.
	Vars map[string]string
	// Changes is the number of entries that were actually modified.
	Changes int
}

// Trim applies whitespace trimming to src according to opts and returns a new
// map. The original map is never mutated. An error is returned if opts requests
// neither keys nor values to be trimmed.
func Trim(src map[string]string, opts Options) (Result, error) {
	if !opts.TrimKeys && !opts.TrimValues {
		return Result{}, errors.New("envtrim: at least one of TrimKeys or TrimValues must be true")
	}

	out := make(map[string]string, len(src))
	changes := 0

	for k, v := range src {
		newKey := k
		newVal := v

		if opts.TrimKeys {
			newKey = strings.TrimSpace(k)
		}
		if opts.TrimValues {
			newVal = strings.TrimSpace(v)
		}

		if newKey != k || newVal != v {
			changes++
		}

		out[newKey] = newVal
	}

	return Result{Vars: out, Changes: changes}, nil
}
