package truncate

import "fmt"

// Options controls truncation behaviour.
type Options struct {
	MaxLen  int
	Suffix  string
	KeepAll bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxLen: 64,
		Suffix: "...",
	}
}

// Result holds a single truncation result.
type Result struct {
	Key       string
	Original  string
	Truncated string
	WasCut    bool
}

// Truncator applies value truncation to env maps.
type Truncator struct {
	opts Options
}

// NewTruncator creates a Truncator with the given options.
func NewTruncator(opts Options) *Truncator {
	if opts.MaxLen <= 0 {
		opts.MaxLen = DefaultOptions().MaxLen
	}
	if opts.Suffix == "" {
		opts.Suffix = DefaultOptions().Suffix
	}
	return &Truncator{opts: opts}
}

// Apply truncates values in the provided map and returns results.
// The original map is not mutated.
func (t *Truncator) Apply(vars map[string]string) (map[string]string, []Result) {
	out := make(map[string]string, len(vars))
	results := make([]Result, 0)
	for k, v := range vars {
		trunc, cut := t.truncate(v)
		out[k] = trunc
		results = append(results, Result{Key: k, Original: v, Truncated: trunc, WasCut: cut})
	}
	return out, results
}

func (t *Truncator) truncate(v string) (string, bool) {
	if t.opts.KeepAll || len(v) <= t.opts.MaxLen {
		return v, false
	}
	cutAt := t.opts.MaxLen - len(t.opts.Suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return v[:cutAt] + t.opts.Suffix, true
}

// Format returns a human-readable summary of truncation results.
func Format(results []Result) string {
	cut := 0
	for _, r := range results {
		if r.WasCut {
			cut++
		}
	}
	return fmt.Sprintf("%d value(s) truncated out of %d", cut, len(results))
}
