package envrank

import (
	"fmt"
	"sort"
)

// Strategy defines how variables are ranked.
type Strategy string

const (
	StrategyKeyLen   Strategy = "key_len"
	StrategyValueLen Strategy = "value_len"
	StrategyAlpha    Strategy = "alpha"
)

var supported = []Strategy{StrategyKeyLen, StrategyValueLen, StrategyAlpha}

// Supported returns all valid ranking strategies.
func Supported() []Strategy { return supported }

// IsSupported reports whether s is a known strategy.
func IsSupported(s Strategy) bool {
	for _, v := range supported {
		if v == s {
			return true
		}
	}
	return false
}

// Entry holds a key/value pair with its computed rank.
type Entry struct {
	Key   string
	Value string
	Rank  int
}

// Rank orders the entries in vars by the given strategy and returns
// a slice of Entry values sorted from highest rank to lowest.
func Rank(vars map[string]string, s Strategy) ([]Entry, error) {
	if !IsSupported(s) {
		return nil, fmt.Errorf("envrank: unsupported strategy %q", s)
	}

	entries := make([]Entry, 0, len(vars))
	for k, v := range vars {
		var rank int
		switch s {
		case StrategyKeyLen:
			rank = len(k)
		case StrategyValueLen:
			rank = len(v)
		case StrategyAlpha:
			// rank by position: lower alphabetical order = higher rank (lower index)
			rank = -alphabeticScore(k)
		}
		entries = append(entries, Entry{Key: k, Value: v, Rank: rank})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rank != entries[j].Rank {
			return entries[i].Rank > entries[j].Rank
		}
		return entries[i].Key < entries[j].Key
	})

	return entries, nil
}

// alphabeticScore returns a negative integer so that 'A' sorts before 'Z'.
func alphabeticScore(key string) int {
	if len(key) == 0 {
		return 0
	}
	// Use the first byte as a simple proxy for alphabetic position.
	return int(key[0])
}
