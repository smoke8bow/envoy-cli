package draft

import "fmt"

// ProfileWriter is the minimal interface needed to persist a profile.
type ProfileWriter interface {
	Save(name string, vars map[string]string) error
}

// Committer wraps a Manager and a ProfileWriter so that a draft can be
// finalised and written to persistent storage in one step.
type Committer struct {
	mgr    *Manager
	writer ProfileWriter
}

// NewCommitter returns a Committer backed by the given Manager and writer.
func NewCommitter(mgr *Manager, writer ProfileWriter) *Committer {
	return &Committer{mgr: mgr, writer: writer}
}

// Commit writes the named draft to the ProfileWriter under profileName and
// then discards the draft. If profileName is empty the draft name is used.
func (c *Committer) Commit(draftName, profileName string) error {
	if profileName == "" {
		profileName = draftName
	}
	vars, err := c.mgr.Get(draftName)
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	if err := c.writer.Save(profileName, vars); err != nil {
		return fmt.Errorf("commit: save profile %q: %w", profileName, err)
	}
	return c.mgr.Discard(draftName)
}
