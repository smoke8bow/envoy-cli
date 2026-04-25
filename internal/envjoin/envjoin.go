// Package envjoin merges multiple env maps into one, with configurable
// separator for values that share the same key (e.g. PATH-style joining).
package envjoin

import (
	"errors"
	"fmt"
	"strings"
)

// Options controls how envjoin behaves.
type Options struct {
	// Separator is placed between values when the same key appears in
	// multiple sources. Defaults to ":".
	Separator string
	// Deduplicate removes duplicate segments after joining.
	Deduplicate bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:   ":",
		Deduplicate: false,
	}
}

// Join merges all sources left-to-right. Keys that appear in more than one
// source are joined with opts.Separator. Keys that appear only once are kept
// as-is.
func Join(opts Options, sources ...map[string]string) (map[string]string, error) {
	if opts.Separator == "" {
		return nil, errors.New("envjoin: separator must not be empty")
	}

	// accumulate segments per key preserving insertion order of keys
	order := []string{}
	segments := map[string][]string{}

	for _, src := range sources {
		for k, v := range src {
			if _, seen := segments[k]; !seen {
				order = append(order, k)
			}
			segments[k] = append(segments[k], v)
		}
	}

	out := make(map[string]string, len(order))
	for _, k := range order {
		parts := segments[k]
		if opts.Deduplicate {
			parts = dedup(parts)
		}
		out[k] = strings.Join(parts, opts.Separator)
	}
	return out, nil
}

// Format returns a human-readable summary of how many keys were joined vs kept.
func Format(result map[string]string, sources []map[string]string) string {
	joined := 0
	for k := range result {
		count := 0
		for _, src := range sources {
			if _, ok := src[k]; ok {
				count++
			}
		}
		if count > 1 {
			joined++
		}
	}
	return fmt.Sprintf("%d total keys, %d joined across sources", len(result), joined)
}

func dedup(parts []string) []string {
	seen := map[string]struct{}{}
	out := parts[:0:0]
	for _, p := range parts {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}
