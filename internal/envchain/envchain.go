// Package envchain provides chained environment variable resolution
// across multiple named profiles, merging them in priority order.
package envchain

import "fmt"

// Link represents a single profile in the chain with its priority.
type Link struct {
	Profile  string
	Vars     map[string]string
	Priority int // higher value = higher priority
}

// Chain holds an ordered list of profile links.
type Chain struct {
	links []Link
}

// ProfileGetter retrieves environment variables for a named profile.
type ProfileGetter interface {
	Get(name string) (map[string]string, error)
}

// NewChain builds a Chain by loading each named profile via getter.
// Profiles listed earlier have lower priority; later entries win on conflict.
func NewChain(getter ProfileGetter, profiles []string) (*Chain, error) {
	if len(profiles) == 0 {
		return nil, fmt.Errorf("envchain: at least one profile required")
	}
	c := &Chain{}
	for i, name := range profiles {
		vars, err := getter.Get(name)
		if err != nil {
			return nil, fmt.Errorf("envchain: load profile %q: %w", name, err)
		}
		c.links = append(c.links, Link{
			Profile:  name,
			Vars:     vars,
			Priority: i,
		})
	}
	return c, nil
}

// Resolve merges all links in priority order and returns the final map.
// Keys from higher-priority profiles overwrite lower-priority ones.
func (c *Chain) Resolve() map[string]string {
	out := make(map[string]string)
	for _, link := range c.links {
		for k, v := range link.Vars {
			out[k] = v
		}
	}
	return out
}

// Source returns which profile a given key originates from (highest priority wins).
// Returns an empty string if the key is not present in any profile.
func (c *Chain) Source(key string) string {
	source := ""
	for _, link := range c.links {
		if _, ok := link.Vars[key]; ok {
			source = link.Profile
		}
	}
	return source
}

// Links returns a copy of the chain's links for inspection.
func (c *Chain) Links() []Link {
	out := make([]Link, len(c.links))
	copy(out, c.links)
	return out
}
