package flatten

import (
	"fmt"
	"strings"
)

// Options controls how nested keys are flattened.
type Options struct {
	Separator string // default "_"
	Uppercase bool   // convert keys to uppercase
	Prefix    string // optional prefix for all keys
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		Uppercase: true,
	}
}

// Flatten takes a nested map[string]any and returns a flat map[string]string.
// Only string leaf values are included; others are formatted with fmt.Sprintf.
func Flatten(nested map[string]any, opts Options) map[string]string {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	result := make(map[string]string)
	flattenRecurse(nested, opts.Prefix, opts, result)
	return result
}

func flattenRecurse(m map[string]any, prefix string, opts Options, out map[string]string) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + opts.Separator + k
		}
		if opts.Uppercase {
			key = strings.ToUpper(key)
		}
		switch val := v.(type) {
		case map[string]any:
			flattenRecurse(val, key, opts, out)
		case string:
			out[key] = val
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
}
