package extract

import "fmt"

// Options controls how keys are extracted from a profile's env map.
type Options struct {
	// Keys is the explicit list of keys to extract. If empty, all keys are returned.
	Keys []string
	// FailOnMissing causes Extract to return an error if any requested key is absent.
	FailOnMissing bool
}

// Result holds the extracted subset of environment variables and any keys
// that were requested but not found in the source map.
type Result struct {
	Vars    map[string]string
	Missing []string
}

// Extract returns a subset of src according to opts.
// When opts.Keys is empty the full map is returned (shallow copy).
func Extract(src map[string]string, opts Options) (Result, error) {
	if len(opts.Keys) == 0 {
		copy := make(map[string]string, len(src))
		for k, v := range src {
			copy[k] = v
		}
		return Result{Vars: copy}, nil
	}

	out := make(map[string]string, len(opts.Keys))
	var missing []string

	for _, k := range opts.Keys {
		v, ok := src[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		out[k] = v
	}

	if opts.FailOnMissing && len(missing) > 0 {
		return Result{}, fmt.Errorf("extract: missing keys: %v", missing)
	}

	return Result{Vars: out, Missing: missing}, nil
}
