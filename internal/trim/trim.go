package trim

import (
	"fmt"
	"strings"
)

// Result holds the outcome of a trim operation.
type Result struct {
	Removed []string
	Kept    map[string]string
}

// Trimmer removes env vars whose values exceed a maximum byte length.
type Trimmer struct {
	maxLen int
}

// NewTrimmer creates a Trimmer with the given max value length.
// If maxLen <= 0 a default of 256 is used.
func NewTrimmer(maxLen int) *Trimmer {
	if maxLen <= 0 {
		maxLen = 256
	}
	return &Trimmer{maxLen: maxLen}
}

// Apply returns a copy of vars with oversized entries removed.
func (t *Trimmer) Apply(vars map[string]string) Result {
	kept := make(map[string]string, len(vars))
	var removed []string
	for k, v := range vars {
		if len(v) > t.maxLen {
			removed = append(removed, k)
		} else {
			kept[k] = v
		}
	}
	return Result{Removed: removed, Kept: kept}
}

// TrimValues truncates values that exceed maxLen instead of removing them.
func (t *Trimmer) TrimValues(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if len(v) > t.maxLen {
			out[k] = v[:t.maxLen]
		} else {
			out[k] = v
		}
	}
	return out
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Removed) == 0 {
		return "no oversized variables found"
	}
	return fmt.Sprintf("removed %d oversized variable(s): %s",
		len(r.Removed), strings.Join(r.Removed, ", "))
}
