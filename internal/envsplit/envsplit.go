// Package envsplit provides utilities for splitting a flat env map
// into multiple named buckets based on key-prefix rules.
package envsplit

import (
	"fmt"
	"strings"
)

// Rule maps a key prefix to a destination bucket name.
type Rule struct {
	Prefix string
	Bucket string
}

// Result holds the output of a Split operation.
type Result struct {
	// Buckets contains the split env maps keyed by bucket name.
	Buckets map[string]map[string]string
	// Remainder holds keys that did not match any rule.
	Remainder map[string]string
}

// Options controls Split behaviour.
type Options struct {
	// StripPrefix removes the matched prefix from keys in the bucket.
	StripPrefix bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{StripPrefix: false}
}

// Split divides src into buckets according to rules.
// Rules are evaluated in order; the first match wins.
// Keys that match no rule are placed in Remainder.
func Split(src map[string]string, rules []Rule, opts Options) (Result, error) {
	for i, r := range rules {
		if r.Prefix == "" {
			return Result{}, fmt.Errorf("rule %d: prefix must not be empty", i)
		}
		if r.Bucket == "" {
			return Result{}, fmt.Errorf("rule %d: bucket must not be empty", i)
		}
	}

	res := Result{
		Buckets:   make(map[string]map[string]string),
		Remainder: make(map[string]string),
	}

	for k, v := range src {
		matched := false
		for _, r := range rules {
			if strings.HasPrefix(k, r.Prefix) {
				if res.Buckets[r.Bucket] == nil {
					res.Buckets[r.Bucket] = make(map[string]string)
				}
				key := k
				if opts.StripPrefix {
					key = strings.TrimPrefix(k, r.Prefix)
				}
				res.Buckets[r.Bucket][key] = v
				matched = true
				break
			}
		}
		if !matched {
			res.Remainder[k] = v
		}
	}

	return res, nil
}
