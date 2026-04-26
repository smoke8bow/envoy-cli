package envlabel

import (
	"fmt"
	"strings"
)

// FormatOptions controls how labels are rendered.
type FormatOptions struct {
	ShowProfile bool
	Separator   string // default "="
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		ShowProfile: false,
		Separator:   "=",
	}
}

// Format renders labels as human-readable lines.
// If ShowProfile is true, each line is prefixed with "<profile>: ".
func Format(profile string, labels []Label, opts FormatOptions) string {
	if len(labels) == 0 {
		return ""
	}
	sep := opts.Separator
	if sep == "" {
		sep = "="
	}
	var sb strings.Builder
	for i, l := range labels {
		if i > 0 {
			sb.WriteByte('\n')
		}
		if opts.ShowProfile {
			fmt.Fprintf(&sb, "%s: %s%s%s", profile, l.Key, sep, l.Value)
		} else {
			fmt.Fprintf(&sb, "%s%s%s", l.Key, sep, l.Value)
		}
	}
	return sb.String()
}
