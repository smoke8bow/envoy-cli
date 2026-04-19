package pivot

import "fmt"

// Direction controls how a profile is pivoted.
type Direction string

const (
	DirectionKeysToValues Direction = "keys-to-values"
	DirectionValuesToKeys Direction = "values-to-keys"
)

// Supported returns all valid pivot directions.
func Supported() []Direction {
	return []Direction{DirectionKeysToValues, DirectionValuesToKeys}
}

// IsSupported returns true if d is a known direction.
func IsSupported(d Direction) bool {
	for _, s := range Supported() {
		if s == d {
			return true
		}
	}
	return false
}

// Pivot transforms a map according to the given direction.
// DirectionKeysToValues swaps each key and value.
// DirectionValuesToKeys is an alias for the same swap operation.
// An error is returned if duplicate values exist (they would collide as keys).
func Pivot(vars map[string]string, dir Direction) (map[string]string, error) {
	if !IsSupported(dir) {
		return nil, fmt.Errorf("pivot: unsupported direction %q", dir)
	}
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if _, exists := out[v]; exists {
			return nil, fmt.Errorf("pivot: duplicate value %q would produce conflicting key", v)
		}
		out[v] = k
	}
	return out, nil
}
