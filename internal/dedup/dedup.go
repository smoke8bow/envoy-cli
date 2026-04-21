// Package dedup provides utilities for detecting and removing duplicate
// environment variable values within a profile.
package dedup

import (
	"fmt"
	"sort"
)

// Strategy controls how duplicates are handled.
type Strategy string

const (
	StrategyKeepFirst Strategy = "keep-first"
	StrategyKeepLast  Strategy = "keep-last"
	StrategyRemoveAll Strategy = "remove-all"
)

// Result describes a group of keys that share the same value.
type Result struct {
	Value string
	Keys  []string
}

// Deduplicator finds and removes keys with duplicate values.
type Deduplicator struct {
	strategy Strategy
}

// Supported returns all valid strategy names.
func Supported() []string {
	return []string{string(StrategyKeepFirst), string(StrategyKeepLast), string(StrategyRemoveAll)}
}

// New returns a Deduplicator with the given strategy.
func New(strategy Strategy) (*Deduplicator, error) {
	switch strategy {
	case StrategyKeepFirst, StrategyKeepLast, StrategyRemoveAll:
		return &Deduplicator{strategy: strategy}, nil
	default:
		return nil, fmt.Errorf("dedup: unsupported strategy %q", strategy)
	}
}

// Find returns all groups of keys that share the same value.
func Find(vars map[string]string) []Result {
	inverted := make(map[string][]string)
	for k, v := range vars {
		inverted[v] = append(inverted[v], k)
	}
	var results []Result
	for v, keys := range inverted {
		if len(keys) > 1 {
			sort.Strings(keys)
			results = append(results, Result{Value: v, Keys: keys})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Value < results[j].Value
	})
	return results
}

// Apply removes duplicate-value keys from vars according to the strategy,
// returning a new map without mutating the input.
func (d *Deduplicator) Apply(vars map[string]string) map[string]string {
	groups := Find(vars)
	remove := make(map[string]bool)
	for _, g := range groups {
		switch d.strategy {
		case StrategyKeepFirst:
			for _, k := range g.Keys[1:] {
				remove[k] = true
			}
		case StrategyKeepLast:
			for _, k := range g.Keys[:len(g.Keys)-1] {
				remove[k] = true
			}
		case StrategyRemoveAll:
			for _, k := range g.Keys {
				remove[k] = true
			}
		}
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if !remove[k] {
			out[k] = v
		}
	}
	return out
}
