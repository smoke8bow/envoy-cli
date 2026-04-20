package diff2

import (
	"fmt"
	"strings"
)

// FormatOptions controls how the diff output is rendered.
type FormatOptions struct {
	LeftLabel  string
	RightLabel string
	ShowEqual  bool
}

// DefaultFormatOptions returns sensible defaults for formatting.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		LeftLabel:  "left",
		RightLabel: "right",
		ShowEqual:  false,
	}
}

// Format renders a Result as a human-readable string.
func Format(r Result, opts FormatOptions) string {
	var sb strings.Builder

	for _, e := range r.Entries {
		switch e.Side {
		case SideLeft:
			sb.WriteString(fmt.Sprintf("- [%s only] %s=%s\n", opts.LeftLabel, e.Key, e.Left))
		case SideRight:
			sb.WriteString(fmt.Sprintf("+ [%s only] %s=%s\n", opts.RightLabel, e.Key, e.Right))
		case SideBoth:
			if e.Equal {
				if opts.ShowEqual {
					sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, e.Left))
				}
			} else {
				sb.WriteString(fmt.Sprintf("~ %s: %s=%s → %s=%s\n",
					e.Key, opts.LeftLabel, e.Left, opts.RightLabel, e.Right))
			}
		}
	}

	return sb.String()
}

// Summary returns a one-line summary of the diff result.
func Summary(r Result) string {
	added := len(r.OnlyRight())
	removed := len(r.OnlyLeft())
	changed := len(r.Changed())

	if added == 0 && removed == 0 && changed == 0 {
		return "no differences"
	}

	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if removed > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", removed))
	}
	if changed > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", changed))
	}
	return strings.Join(parts, ", ")
}
