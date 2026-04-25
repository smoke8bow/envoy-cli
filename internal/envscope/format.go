package envscope

import (
	"fmt"
	"strings"
)

// FormatOptions controls how a ScopedView is rendered.
type FormatOptions struct {
	// ShowProfile includes a header line with the profile and scope names.
	ShowProfile bool
}

// DefaultFormatOptions returns sensible defaults.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{ShowProfile: true}
}

// Format renders a ScopedView as a human-readable string.
// Keys are printed in sorted order.
func Format(sv *ScopedView, opts FormatOptions) string {
	var sb strings.Builder

	if opts.ShowProfile {
		fmt.Fprintf(&sb, "# profile=%s scope=%s\n", sv.Profile, sv.Scope)
	}

	for _, k := range sv.Keys() {
		fmt.Fprintf(&sb, "%s=%s\n", k, sv.Vars[k])
	}

	return sb.String()
}
