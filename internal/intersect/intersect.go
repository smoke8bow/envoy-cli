// Package intersect computes the intersection of two environment variable maps,
// returning only the keys that exist in both profiles.
package intersect

import "sort"

// Result holds the outcome of an intersection operation.
type Result struct {
	// Vars contains keys present in both A and B, with values taken from A.
	Vars map[string]string
	// Conflicts contains keys where both maps have different values.
	Conflicts map[string]Conflict
}

// Conflict describes a key whose value differs between the two profiles.
type Conflict struct {
	ValueA string
	ValueB string
}

// Intersect returns keys common to both a and b.
// Values in the result are taken from a.
// Keys whose values differ are recorded in Conflicts.
func Intersect(a, b map[string]string) Result {
	vars := make(map[string]string)
	conflicts := make(map[string]Conflict)

	for k, va := range a {
		vb, ok := b[k]
		if !ok {
			continue
		}
		vars[k] = va
		if va != vb {
			conflicts[k] = Conflict{ValueA: va, ValueB: vb}
		}
	}

	return Result{Vars: vars, Conflicts: conflicts}
}

// Keys returns a sorted slice of keys in the result.
func (r Result) Keys() []string {
	keys := make([]string, 0, len(r.Vars))
	for k := range r.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ConflictKeys returns a sorted slice of keys that have conflicting values.
func (r Result) ConflictKeys() []string {
	keys := make([]string, 0, len(r.Conflicts))
	for k := range r.Conflicts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
