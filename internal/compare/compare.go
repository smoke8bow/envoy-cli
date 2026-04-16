package compare

import (
	"fmt"
	"sort"
)

// Result holds the comparison between two profiles.
type Result struct {
	OnlyInA  map[string]string
	OnlyInB  map[string]string
	Different map[string]Pair
	Same     map[string]string
}

// Pair holds the two differing values for a key.
type Pair struct {
	A string
	B string
}

// Compare returns a Result describing differences between two env maps.
func Compare(a, b map[string]string) Result {
	r := Result{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Different: make(map[string]Pair),
		Same:      make(map[string]string),
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok {
			r.OnlyInA[k] = v
		} else if v != bv {
			r.Different[k] = Pair{A: v, B: bv}
		} else {
			r.Same[k] = v
		}
	}
	for k, v := range b {
		if _, ok := a[k]; !ok {
			r.OnlyInB[k] = v
		}
	}
	return r
}

// Format returns a human-readable summary of the Result.
func Format(nameA, nameB string, r Result) string {
	out := fmt.Sprintf("Comparing [%s] vs [%s]\n", nameA, nameB)

	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m { ks = append(ks, k) }
		sort.Strings(ks)
		return ks
	}

	for _, k := range keys(r.OnlyInA) {
		out += fmt.Sprintf("  < %s=%s\n", k, r.OnlyInA[k])
	}
	for _, k := range keys(r.OnlyInB) {
		out += fmt.Sprintf("  > %s=%s\n", k, r.OnlyInB[k])
	}
	dks := make([]string, 0, len(r.Different))
	for k := range r.Different { dks = append(dks, k) }
	sort.Strings(dks)
	for _, k := range dks {
		p := r.Different[k]
		out += fmt.Sprintf("  ~ %s: %s -> %s\n", k, p.A, p.B)
	}
	return out
}
