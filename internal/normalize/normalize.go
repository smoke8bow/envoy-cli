package normalize

import (
	"errors"
	"strings"
)

// Strategy defines how keys should be normalized.
type Strategy string

const (
	StrategyUpper  Strategy = "upper"
	StrategyLower  Strategy = "lower"
	StrategySnake  Strategy = "snake"
)

var supported = []Strategy{StrategyUpper, StrategyLower, StrategySnake}

func Supported() []Strategy { return supported }

func IsSupported(s Strategy) bool {
	for _, v := range supported {
		if v == s {
			return true
		}
	}
	return false
}

// Normalizer applies a key normalization strategy to an env map.
type Normalizer struct {
	strategy Strategy
}

func New(s Strategy) (*Normalizer, error) {
	if !IsSupported(s) {
		return nil, errors.New("normalize: unsupported strategy: " + string(s))
	}
	return &Normalizer{strategy: s}, nil
}

// Apply returns a new map with keys normalized according to the strategy.
// If two keys collide after normalization, the last one wins.
func (n *Normalizer) Apply(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[n.normalizeKey(k)] = v
	}
	return out
}

func (n *Normalizer) normalizeKey(k string) string {
	switch n.strategy {
	case StrategyUpper:
		return strings.ToUpper(k)
	case StrategyLower:
		return strings.ToLower(k)
	case StrategySnake:
		return toSnake(k)
	}
	return k
}

func toSnake(k string) string {
	k = strings.ToUpper(k)
	k = strings.ReplaceAll(k, "-", "_")
	k = strings.ReplaceAll(k, " ", "_")
	return k
}
