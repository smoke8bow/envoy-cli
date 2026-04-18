package promote

import "fmt"

// Store describes the subset of store operations needed by the Promoter.
type Store interface {
	Get(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
	List() ([]string, error)
}

// Promoter copies selected keys from a source profile into a destination
// profile, optionally overwriting existing keys.
type Promoter struct {
	store Store
}

// NewPromoter returns a Promoter backed by the given store.
func NewPromoter(s Store) *Promoter {
	return &Promoter{store: s}
}

// PromoteOptions controls how promotion behaves.
type PromoteOptions struct {
	// Keys is the list of keys to promote. If empty, all keys are promoted.
	Keys []string
	// Overwrite controls whether existing keys in dst are replaced.
	Overwrite bool
}

// Promote copies keys from src into dst according to opts.
// Returns the list of keys that were actually written.
func (p *Promoter) Promote(src, dst string, opts PromoteOptions) ([]string, error) {
	srcVars, err := p.store.Get(src)
	if err != nil {
		return nil, fmt.Errorf("promote: source %q: %w", src, err)
	}
	dstVars, err := p.store.Get(dst)
	if err != nil {
		return nil, fmt.Errorf("promote: destination %q: %w", dst, err)
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range srcVars {
			keys = append(keys, k)
		}
	}

	var written []string
	for _, k := range keys {
		v, ok := srcVars[k]
		if !ok {
			return nil, fmt.Errorf("promote: key %q not found in source %q", k, src)
		}
		if _, exists := dstVars[k]; exists && !opts.Overwrite {
			continue
		}
		dstVars[k] = v
		written = append(written, k)
	}

	if err := p.store.Save(dst, dstVars); err != nil {
		return nil, fmt.Errorf("promote: save %q: %w", dst, err)
	}
	return written, nil
}
