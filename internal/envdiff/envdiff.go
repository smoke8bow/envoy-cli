// Package envdiff compares two profiles against the current OS environment.
package envdiff

import (
	"fmt"
	"os"
	"sort"
)

// Source indicates where a variable originates.
type Source string

const (
	SourceProfile Source = "profile"
	SourceOS      Source = "os"
	SourceBoth    Source = "both"
)

// Entry represents a single variable in the diff result.
type Entry struct {
	Key          string
	ProfileValue string
	OSValue      string
	Source       Source
}

// Result holds the outcome of comparing a profile against the OS environment.
type Result struct {
	OnlyInProfile []Entry
	OnlyInOS      []Entry
	InBoth        []Entry
}

// Compare contrasts profileVars against the current OS environment.
// Keys present in both are recorded with both values for inspection.
func Compare(profileVars map[string]string) Result {
	osVars := osEnvMap()

	var result Result

	for k, pv := range profileVars {
		if ov, ok := osVars[k]; ok {
			result.InBoth = append(result.InBoth, Entry{
				Key: k, ProfileValue: pv, OSValue: ov, Source: SourceBoth,
			})
		} else {
			result.OnlyInProfile = append(result.OnlyInProfile, Entry{
				Key: k, ProfileValue: pv, Source: SourceProfile,
			})
		}
	}

	for k, ov := range osVars {
		if _, ok := profileVars[k]; !ok {
			result.OnlyInOS = append(result.OnlyInOS, Entry{
				Key: k, OSValue: ov, Source: SourceOS,
			})
		}
	}

	sortEntries(result.OnlyInProfile)
	sortEntries(result.OnlyInOS)
	sortEntries(result.InBoth)

	return result
}

// Format renders a human-readable summary of the diff result.
func Format(r Result) string {
	out := ""
	for _, e := range r.OnlyInProfile {
		out += fmt.Sprintf("+ %s=%s\n", e.Key, e.ProfileValue)
	}
	for _, e := range r.InBoth {
		if e.ProfileValue != e.OSValue {
			out += fmt.Sprintf("~ %s: profile=%q os=%q\n", e.Key, e.ProfileValue, e.OSValue)
		} else {
			out += fmt.Sprintf("= %s=%s\n", e.Key, e.ProfileValue)
		}
	}
	for _, e := range r.OnlyInOS {
		out += fmt.Sprintf("- %s=%s\n", e.Key, e.OSValue)
	}
	return out
}

func osEnvMap() map[string]string {
	m := make(map[string]string)
	for _, pair := range os.Environ() {
		for i := 0; i < len(pair); i++ {
			if pair[i] == '=' {
				m[pair[:i]] = pair[i+1:]
				break
			}
		}
	}
	return m
}

func sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
}
