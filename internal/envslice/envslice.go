// Package envslice provides utilities for converting environment variable
// maps to ordered slices and back, with support for filtering and sorting.
package envslice

import (
	"fmt"
	"sort"
	"strings"
)

// Entry represents a single environment variable as a key=value pair.
type Entry struct {
	Key   string
	Value string
}

// String returns the entry in KEY=VALUE format.
func (e Entry) String() string {
	return fmt.Sprintf("%s=%s", e.Key, e.Value)
}

// FromMap converts a map of environment variables to a sorted slice of Entry.
func FromMap(vars map[string]string) []Entry {
	entries := make([]Entry, 0, len(vars))
	for k, v := range vars {
		entries = append(entries, Entry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// ToMap converts a slice of Entry back to a map.
// Later entries with duplicate keys overwrite earlier ones.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// ToStrings converts a slice of Entry to KEY=VALUE strings.
func ToStrings(entries []Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.String()
	}
	return out
}

// FromStrings parses a slice of KEY=VALUE strings into entries.
// Lines without '=' are skipped.
func FromStrings(lines []string) []Entry {
	var entries []Entry
	for _, line := range lines {
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		entries = append(entries, Entry{
			Key:   line[:idx],
			Value: line[idx+1:],
		})
	}
	return entries
}

// FilterByPrefix returns only entries whose key starts with the given prefix.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	var out []Entry
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			out = append(out, e)
		}
	}
	return out
}
