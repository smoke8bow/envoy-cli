package clone

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/profile"
	"github.com/envoy-cli/envoy/internal/validate"
)

// Cloner handles duplicating profiles under a new name.
type Cloner struct {
	manager *profile.Manager
}

// NewCloner returns a Cloner backed by the given profile manager.
func NewCloner(m *profile.Manager) *Cloner {
	return &Cloner{manager: m}
}

// Clone copies the vars from src into a new profile named dst.
// It returns an error if src does not exist, dst already exists,
// or dst is not a valid profile name.
func (c *Cloner) Clone(src, dst string) error {
	if err := validate.ProfileName(dst); err != nil {
		return fmt.Errorf("invalid destination name: %w", err)
	}

	srcProfile, err := c.manager.Get(src)
	if err != nil {
		return fmt.Errorf("source profile %q not found: %w", src, err)
	}

	// Copy vars so mutations to the clone don't affect the original.
	copied := make(map[string]string, len(srcProfile.Vars))
	for k, v := range srcProfile.Vars {
		copied[k] = v
	}

	if err := c.manager.Create(dst, copied); err != nil {
		return fmt.Errorf("create clone %q: %w", dst, err)
	}

	return nil
}

// CloneWithOverrides copies src into dst and then merges overrides
// on top of the copied vars before saving.
func (c *Cloner) CloneWithOverrides(src, dst string, overrides map[string]string) error {
	if err := c.Clone(src, dst); err != nil {
		return err
	}

	if len(overrides) == 0 {
		return nil
	}

	cloned, err := c.manager.Get(dst)
	if err != nil {
		return err
	}

	merged := make(map[string]string, len(cloned.Vars)+len(overrides))
	for k, v := range cloned.Vars {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}

	return c.manager.Update(dst, merged)
}
