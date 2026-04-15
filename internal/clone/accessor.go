package clone

import "github.com/envoy-cli/envoy/internal/profile"

// Manager exposes the underlying profile.Manager so tests and callers
// can inspect state without breaking encapsulation of Cloner.
func (c *Cloner) Manager() *profile.Manager {
	return c.manager
}
