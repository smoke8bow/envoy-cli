package mask

import "strings"

// DefaultPatterns are key substrings that trigger masking.
var DefaultPatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE"}

// Masker redacts sensitive environment variable values.
type Masker struct {
	patterns []string
}

// NewMasker returns a Masker using the provided patterns.
// If patterns is nil, DefaultPatterns are used.
func NewMasker(patterns []string) *Masker {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return &Masker{patterns: upper}
}

// IsSensitive reports whether the key matches any pattern.
func (m *Masker) IsSensitive(key string) bool {
	k := strings.ToUpper(key)
	for _, p := range m.patterns {
		if strings.Contains(k, p) {
			return true
		}
	}
	return false
}

// MaskValue returns "***" if the key is sensitive, otherwise the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return "***"
	}
	return value
}

// Apply returns a copy of vars with sensitive values replaced.
func (m *Masker) Apply(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = m.MaskValue(k, v)
	}
	return out
}
