package sort

import (
	"errors"
	"sort"
)

// Strategy defines how profile keys should be sorted.
type Strategy string

const (
	StrategyAlpha  Strategy = "alpha"
	StrategyReverse Strategy = "reverse"
	StrategyLength  Strategy = "length"
)

var ErrUnknownStrategy = errors.New("unknown sort strategy")

// Supported returns all valid strategies.
func Supported() []Strategy {
	return []Strategy{StrategyAlpha, StrategyReverse, StrategyLength}
}

// IsSupported reports whether s is a known strategy.
func IsSupported(s Strategy) bool {
	for _, v := range Supported() {
		if v == s {
			return true
		}
	}
	return false
}

// Apply returns a new map with the same key/value pairs; the returned
// slice of keys reflects the requested ordering.
func Apply(vars map[string]string, s Strategy) ([]string, error) {
	if !IsSupported(s) {
		return nil, ErrUnknownStrategy
	}
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	switch s {
	case StrategyAlpha:
		sort.Strings(keys)
	case StrategyReverse:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case StrategyLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	}
	return keys, nil
}
