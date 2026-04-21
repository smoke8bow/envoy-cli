package envstats

import (
	"fmt"
	"sort"
	"strings"
)

// Stats holds computed statistics for a set of environment variables.
type Stats struct {
	Total       int
	Empty       int
	NonEmpty    int
	AvgKeyLen   float64
	AvgValueLen float64
	LongestKey  string
	ShortestKey string
	LongestVal  string
}

// Compute calculates statistics from the provided env map.
func Compute(env map[string]string) Stats {
	if len(env) == 0 {
		return Stats{}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var totalKeyLen, totalValLen int
	var longestKey, shortestKey, longestVal string
	empty := 0

	for i, k := range keys {
		v := env[k]
		totalKeyLen += len(k)
		totalValLen += len(v)

		if v == "" {
			empty++
		}

		if i == 0 {
			longestKey = k
			shortestKey = k
			longestVal = v
		} else {
			if len(k) > len(longestKey) {
				longestKey = k
			}
			if len(k) < len(shortestKey) {
				shortestKey = k
			}
			if len(v) > len(longestVal) {
				longestVal = v
			}
		}
	}

	n := len(keys)
	return Stats{
		Total:       n,
		Empty:       empty,
		NonEmpty:    n - empty,
		AvgKeyLen:   float64(totalKeyLen) / float64(n),
		AvgValueLen: float64(totalValLen) / float64(n),
		LongestKey:  longestKey,
		ShortestKey: shortestKey,
		LongestVal:  longestVal,
	}
}

// Format returns a human-readable summary of the stats.
func Format(s Stats) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total vars   : %d\n", s.Total)
	fmt.Fprintf(&sb, "Non-empty    : %d\n", s.NonEmpty)
	fmt.Fprintf(&sb, "Empty        : %d\n", s.Empty)
	fmt.Fprintf(&sb, "Avg key len  : %.1f\n", s.AvgKeyLen)
	fmt.Fprintf(&sb, "Avg val len  : %.1f\n", s.AvgValueLen)
	fmt.Fprintf(&sb, "Longest key  : %s\n", s.LongestKey)
	fmt.Fprintf(&sb, "Shortest key : %s\n", s.ShortestKey)
	fmt.Fprintf(&sb, "Longest val  : %s\n", s.LongestVal)
	return sb.String()
}
