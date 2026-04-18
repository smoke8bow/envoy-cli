package redact

import (
	"regexp"
	"strings"
)

// Rule defines a pattern and replacement for redaction.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies redaction rules to env var maps and strings.
type Redactor struct {
	rules []Rule
}

var defaultPatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api[_]?key`,
	`(?i)private[_]?key`,
}

const defaultReplacement = "[REDACTED]"

// New creates a Redactor with the default sensitive key patterns.
func New() *Redactor {
	rules := make([]Rule, 0, len(defaultPatterns))
	for _, p := range defaultPatterns {
		rules = append(rules, Rule{
			Pattern:     regexp.MustCompile(p),
			Replacement: defaultReplacement,
		})
	}
	return &Redactor{rules: rules}
}

// NewWithPatterns creates a Redactor with custom key patterns.
func NewWithPatterns(patterns []string) (*Redactor, error) {
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Pattern: re, Replacement: defaultReplacement})
	}
	return &Redactor{rules: rules}, nil
}

// IsSensitive reports whether the given key matches any redaction rule.
func (r *Redactor) IsSensitive(key string) bool {
	for _, rule := range r.rules {
		if rule.Pattern.MatchString(key) {
			return true
		}
	}
	return false
}

// Apply returns a copy of vars with sensitive values replaced.
func (r *Redactor) Apply(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if r.IsSensitive(k) {
			out[k] = defaultReplacement
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactString replaces any value in s that looks like a known sensitive value.
func (r *Redactor) RedactString(s string, vars map[string]string) string {
	for k, v := range vars {
		if r.IsSensitive(k) && v != "" && strings.Contains(s, v) {
			s = strings.ReplaceAll(s, v, defaultReplacement)
		}
	}
	return s
}
