// Package envrewrite provides key and value rewriting for environment variable maps
// using regex-based find-and-replace rules.
package envrewrite

import (
	"fmt"
	"regexp"
)

// Rule defines a single rewrite operation applied to env var keys or values.
type Rule struct {
	// Target specifies what to rewrite: "key", "value", or "both".
	Target string
	// Pattern is the regular expression to match against.
	Pattern string
	// Replacement is the replacement string (supports $1, $2 capture groups).
	Replacement string

	compiled *regexp.Regexp
}

// compile prepares the rule's regex, returning an error if the pattern is invalid.
func (r *Rule) compile() error {
	if r.compiled != nil {
		return nil
	}
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("envrewrite: invalid pattern %q: %w", r.Pattern, err)
	}
	r.compiled = re
	return nil
}

// Rewriter applies a set of rewrite rules to an environment variable map.
type Rewriter struct {
	rules []Rule
}

// New creates a Rewriter from the provided rules, compiling all patterns
// upfront. Returns an error if any pattern fails to compile.
func New(rules []Rule) (*Rewriter, error) {
	for i := range rules {
		switch rules[i].Target {
		case "key", "value", "both":
		default:
			return nil, fmt.Errorf("envrewrite: unknown target %q (want key, value, or both)", rules[i].Target)
		}
		if err := rules[i].compile(); err != nil {
			return nil, err
		}
	}
	return &Rewriter{rules: rules}, nil
}

// Apply runs all rewrite rules against src and returns a new map with the
// results. The original map is never mutated.
func (rw *Rewriter) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))

	// Copy first so we work on a stable snapshot.
	for k, v := range src {
		out[k] = v
	}

	for _, rule := range rw.rules {
		rewritten := make(map[string]string, len(out))
		for k, v := range out {
			newKey := k
			newVal := v

			if rule.Target == "key" || rule.Target == "both" {
				newKey = rule.compiled.ReplaceAllString(k, rule.Replacement)
			}
			if rule.Target == "value" || rule.Target == "both" {
				newVal = rule.compiled.ReplaceAllString(v, rule.Replacement)
			}
			rewritten[newKey] = newVal
		}
		out = rewritten
	}
	return out
}

// Summary describes how many keys and values were changed by a rewrite pass.
type Summary struct {
	KeysChanged   int
	ValuesChanged int
}

// Diff compares src against the result of Apply(src) and returns a Summary
// of what changed without retaining the rewritten map.
func (rw *Rewriter) Diff(src map[string]string) Summary {
	result := rw.Apply(src)
	var s Summary
	for k, v := range src {
		newKey := findNewKey(result, k, v)
		if newKey != k {
			s.KeysChanged++
		}
		if nv, ok := result[newKey]; ok && nv != v {
			s.ValuesChanged++
		}
	}
	return s
}

// findNewKey attempts to locate the post-rewrite key that corresponds to the
// original key k with value v. It first checks if k is still present (value
// rewrite only), then falls back to a value scan.
func findNewKey(result map[string]string, k, v string) string {
	if _, ok := result[k]; ok {
		return k
	}
	// Key was renamed; find by matching original value as a heuristic.
	for rk, rv := range result {
		if rv == v {
			_ = rk
			return rk
		}
	}
	return k
}
