package snapshot

// Accessor defines the read-only interface used by other packages
// that need to inspect snapshots without importing the full Manager.
type Accessor interface {
	// Take creates a new snapshot of the named profile.
	Take(profile, note string) (*Entry, error)
	// Restore applies a snapshot's vars back to the target profile.
	Restore(snapshotID, targetProfile string) error
}

// Ensure *Manager satisfies Accessor at compile time.
var _ Accessor = (*Manager)(nil)
