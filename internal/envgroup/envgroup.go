// Package envgroup provides grouping of environment variables within a profile
// by a shared key prefix, enabling structured access and rendering.
package envgroup

import (
	"fmt"
	"sort"
	"strings"
)

// Group holds variables that share a common prefix.
type Group struct {
	Prefix string
	Vars   map[string]string
}

// Result is the output of grouping a flat env map.
type Result struct {
	Groups    []Group
	Ungrouped map[string]string
}

// Options controls grouping behaviour.
type Options struct {
	// MinSize is the minimum number of keys required to form a group.
	// Defaults to 2.
	MinSize int
	// Separator is the delimiter used to detect prefixes (default "_").
	Separator string
	// StripPrefix removes the prefix from keys inside a group when true.
	StripPrefix bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{MinSize: 2, Separator: "_", StripPrefix: false}
}

// GroupBy partitions vars into groups whose keys share a common first segment.
func GroupBy(vars map[string]string, opts Options) Result {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	if opts.MinSize < 1 {
		opts.MinSize = 2
	}

	buckets := map[string]map[string]string{}
	for k, v := range vars {
		parts := strings.SplitN(k, opts.Separator, 2)
		if len(parts) < 2 {
			continue
		}
		prefix := parts[0]
		if buckets[prefix] == nil {
			buckets[prefix] = map[string]string{}
		}
		buckets[prefix][k] = v
	}

	grouped := map[string]bool{}
	var groups []Group
	for prefix, members := range buckets {
		if len(members) < opts.MinSize {
			continue
		}
		out := map[string]string{}
		for k, v := range members {
			key := k
			if opts.StripPrefix {
				key = strings.TrimPrefix(k, prefix+opts.Separator)
			}
			out[key] = v
			grouped[k] = true
		}
		groups = append(groups, Group{Prefix: prefix, Vars: out})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Prefix < groups[j].Prefix })

	ungrouped := map[string]string{}
	for k, v := range vars {
		if !grouped[k] {
			ungrouped[k] = v
		}
	}
	return Result{Groups: groups, Ungrouped: ungrouped}
}

// Format renders a Result as a human-readable string.
func Format(r Result) string {
	var sb strings.Builder
	for _, g := range r.Groups {
		fmt.Fprintf(&sb, "[%s]\n", g.Prefix)
		keys := make([]string, 0, len(g.Vars))
		for k := range g.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, g.Vars[k])
		}
	}
	if len(r.Ungrouped) > 0 {
		sb.WriteString("[ungrouped]\n")
		keys := make([]string, 0, len(r.Ungrouped))
		for k := range r.Ungrouped {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, r.Ungrouped[k])
		}
	}
	return sb.String()
}
