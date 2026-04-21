// Package envcount provides utilities for counting and summarising
// environment variable statistics across one or more profiles.
package envcount

import (
	"fmt"
	"sort"
)

// ProfileStats holds variable counts for a single profile.
type ProfileStats struct {
	Name     string
	Total    int
	Empty    int
	NonEmpty int
}

// String returns a human-readable summary line.
func (s ProfileStats) String() string {
	return fmt.Sprintf("%s: total=%d non-empty=%d empty=%d",
		s.Name, s.Total, s.NonEmpty, s.Empty)
}

// Counter computes statistics for environment variable maps.
type Counter struct{}

// New returns a new Counter.
func New() *Counter { return &Counter{} }

// Compute derives ProfileStats from a named env map.
func (c *Counter) Compute(name string, vars map[string]string) ProfileStats {
	stats := ProfileStats{Name: name, Total: len(vars)}
	for _, v := range vars {
		if v == "" {
			stats.Empty++
		} else {
			stats.NonEmpty++
		}
	}
	return stats
}

// ComputeAll derives stats for multiple profiles and returns them sorted
// alphabetically by profile name.
func (c *Counter) ComputeAll(profiles map[string]map[string]string) []ProfileStats {
	results := make([]ProfileStats, 0, len(profiles))
	for name, vars := range profiles {
		results = append(results, c.Compute(name, vars))
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results
}

// Totals sums stats across all provided ProfileStats entries.
func Totals(stats []ProfileStats) ProfileStats {
	out := ProfileStats{Name: "(total)"}
	for _, s := range stats {
		out.Total += s.Total
		out.Empty += s.Empty
		out.NonEmpty += s.NonEmpty
	}
	return out
}
