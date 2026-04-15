package copy

import (
	"errors"
	"fmt"
)

// Accessor defines the interface for reading and writing profiles.
type Accessor interface {
	Get(name string) (map[string]string, error)
	Create(name string, vars map[string]string) error
	Exists(name string) bool
}

// Copier copies specific keys from one profile into another.
type Copier struct {
	store Accessor
}

// NewCopier returns a new Copier backed by the given store.
func NewCopier(s Accessor) *Copier {
	return &Copier{store: s}
}

// CopyKeys copies the specified keys from src profile into dst profile.
// If dst does not exist it is created. Existing keys in dst are preserved
// unless overwrite is true.
func (c *Copier) CopyKeys(src, dst string, keys []string, overwrite bool) error {
	if src == "" || dst == "" {
		return errors.New("copy: src and dst profile names must not be empty")
	}
	if len(keys) == 0 {
		return errors.New("copy: at least one key must be specified")
	}

	srcVars, err := c.store.Get(src)
	if err != nil {
		return fmt.Errorf("copy: source profile %q not found: %w", src, err)
	}

	var dstVars map[string]string
	if c.store.Exists(dst) {
		dstVars, err = c.store.Get(dst)
		if err != nil {
			return fmt.Errorf("copy: destination profile %q: %w", dst, err)
		}
	} else {
		dstVars = make(map[string]string)
	}

	for _, k := range keys {
		v, ok := srcVars[k]
		if !ok {
			return fmt.Errorf("copy: key %q not found in source profile %q", k, src)
		}
		if _, exists := dstVars[k]; exists && !overwrite {
			continue
		}
		dstVars[k] = v
	}

	if err := c.store.Create(dst, dstVars); err != nil {
		return fmt.Errorf("copy: failed to save destination profile %q: %w", dst, err)
	}
	return nil
}
