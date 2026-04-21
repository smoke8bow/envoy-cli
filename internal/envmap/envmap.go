// Package envmap provides utilities for converting between environment
// variable maps and other representations such as slices and key=value strings.
package envmap

import (
	"fmt"
	"sort"
	"strings"
)

// FromSlice converts a slice of "KEY=VALUE" strings into a map.
// Entries without an "=" are stored with an empty string value.
func FromSlice(pairs []string) map[string]string {
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			out[p] = ""
			continue
		}
		out[p[:idx]] = p[idx+1:]
	}
	return out
}

// ToSlice converts a map into a sorted slice of "KEY=VALUE" strings.
func ToSlice(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(m))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return out
}

// Keys returns a sorted slice of all keys in the map.
func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Clone returns a shallow copy of the map.
func Clone(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Merge merges src into dst, overwriting existing keys. dst is mutated.
func Merge(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}
