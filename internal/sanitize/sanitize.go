// Package sanitize provides utilities for cleaning environment variable
// values by trimming whitespace, removing control characters, and
// optionally enforcing max-length constraints.
package sanitize

import (
	"strings"
	"unicode"
)

// Options controls which sanitisation steps are applied.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// StripControl removes non-printable control characters from values.
	StripControl bool
	// MaxValueLen truncates values longer than this length. 0 means no limit.
	MaxValueLen int
}

// DefaultOptions returns a sensible default sanitisation configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:    true,
		StripControl: true,
		MaxValueLen:  0,
	}
}

// Sanitizer applies sanitisation rules to a map of environment variables.
type Sanitizer struct {
	opts Options
}

// New creates a Sanitizer with the given options.
func New(opts Options) *Sanitizer {
	return &Sanitizer{opts: opts}
}

// Apply returns a new map with sanitised values. The source map is not mutated.
func (s *Sanitizer) Apply(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = s.sanitize(v)
	}
	return out
}

// Value sanitises a single value according to the configured options.
func (s *Sanitizer) Value(v string) string {
	return s.sanitize(v)
}

func (s *Sanitizer) sanitize(v string) string {
	if s.opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if s.opts.StripControl {
		v = stripControl(v)
	}
	if s.opts.MaxValueLen > 0 && len(v) > s.opts.MaxValueLen {
		v = v[:s.opts.MaxValueLen]
	}
	return v
}

func stripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}
		return r
	}, s)
}
