// Package chain provides a mechanism to compose multiple env var maps
// by applying them in sequence, with later entries taking precedence.
package chain

import "fmt"

// Entry is a named layer in the chain.
type Entry struct {
	Name string
	Vars map[string]string
}

// Result holds the merged output and a source map indicating which
// layer each key was resolved from.
type Result struct {
	Vars   map[string]string
	Source map[string]string // key -> layer name
}

// Composer applies a sequence of named env layers in order.
type Composer struct {
	entries []Entry
}

// NewComposer returns a Composer with the provided entries.
// Entries are applied left-to-right; later entries override earlier ones.
func NewComposer(entries []Entry) (*Composer, error) {
	for i, e := range entries {
		if e.Name == "" {
			return nil, fmt.Errorf("chain: entry %d has empty name", i)
		}
		if e.Vars == nil {
			return nil, fmt.Errorf("chain: entry %q has nil vars", e.Name)
		}
	}
	return &Composer{entries: entries}, nil
}

// Compose merges all layers and returns the Result.
func (c *Composer) Compose() Result {
	out := make(map[string]string)
	src := make(map[string]string)
	for _, e := range c.entries {
		for k, v := range e.Vars {
			out[k] = v
			src[k] = e.Name
		}
	}
	return Result{Vars: out, Source: src}
}

// Layers returns the entry names in application order.
func (c *Composer) Layers() []string {
	names := make([]string, len(c.entries))
	for i, e := range c.entries {
		names[i] = e.Name
	}
	return names
}
