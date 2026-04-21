// Package envset provides set-theoretic operations (union, intersection, difference)
// on environment variable maps.
package envset

// Op represents a set operation.
type Op string

const (
	OpUnion        Op = "union"
	OpIntersection Op = "intersection"
	OpDifference   Op = "difference"
)

var supported = []Op{OpUnion, OpIntersection, OpDifference}

// Supported returns all valid operation names.
func Supported() []Op { return supported }

// IsSupported reports whether op is a known operation.
func IsSupported(op Op) bool {
	for _, s := range supported {
		if s == op {
			return true
		}
	}
	return false
}

// Union returns a new map containing all keys from a and b.
// Keys in b override keys in a when both are present.
func Union(a, b map[string]string) map[string]string {
	out := make(map[string]string, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

// Intersection returns a new map containing only keys present in both a and b.
// Values are taken from b.
func Intersection(a, b map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range b {
		if _, ok := a[k]; ok {
			out[k] = v
		}
	}
	return out
}

// Difference returns a new map containing keys that are in a but not in b.
func Difference(a, b map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range a {
		if _, ok := b[k]; !ok {
			out[k] = v
		}
	}
	return out
}

// Apply performs the given operation on a and b.
func Apply(op Op, a, b map[string]string) (map[string]string, error) {
	switch op {
	case OpUnion:
		return Union(a, b), nil
	case OpIntersection:
		return Intersection(a, b), nil
	case OpDifference:
		return Difference(a, b), nil
	default:
		return nil, fmt.Errorf("envset: unsupported operation %q", op)
	}
}
