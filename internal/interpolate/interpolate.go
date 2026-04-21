// Package interpolate provides value interpolation across profile variables,
// supporting cross-profile references and ambient OS environment fallback.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// Getter retrieves a named profile's variables.
type Getter interface {
	Get(name string) (map[string]string, error)
}

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows falling back to os.Getenv when a key is unresolved.
	FallbackToOS bool
	// MaxDepth limits recursive interpolation passes.
	MaxDepth int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS: true,
		MaxDepth:     8,
	}
}

// Interpolator performs variable interpolation within a profile's values.
type Interpolator struct {
	getter  Getter
	opts    Options
	refRe   *regexp.Regexp
}

// New creates an Interpolator with the given getter and options.
func New(g Getter, opts Options) *Interpolator {
	return &Interpolator{
		getter: g,
		opts:   opts,
		// matches ${VAR} and $VAR
		refRe: regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)},
}

// Apply resolves all variable references within vars using vars itself as the
// primary source and the getter for cross-profile lookups.
func (i *Interpolator) Apply(vars map[string]string, ambient func(string) string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	for pass := 0; pass < i.opts.MaxDepth; pass++ {
		changed := false
		for k, v := range out {
			resolved, err := i.resolveValue(v, out, ambient)
			if err != nil {
				return nil, fmt.Errorf("interpolate key %q: %w", k, err)
			}
			if resolved != v {
				out[k] = resolved
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return out, nil
}

func (i *Interpolator) resolveValue(v string, local map[string]string, ambient func(string) string) (string, error) {
	var lastErr error
	result := i.refRe.ReplaceAllStringFunc(v, func(match string) string {
		key := strings.TrimPrefix(strings.TrimPrefix(strings.Trim(match, "${}"), "${"), "$")
		key = strings.TrimSuffix(key, "}")
		if val, ok := local[key]; ok {
			return val
		}
		if ambient != nil && i.opts.FallbackToOS {
			if val := ambient(key); val != "" {
				return val
			}
		}
		lastErr = fmt.Errorf("unresolved variable %q", key)
		return match
	})
	return result, lastErr
}
