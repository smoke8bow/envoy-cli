package rename

import (
	"errors"
	"fmt"
)

// ErrNotFound is returned when the source profile does not exist.
var ErrNotFound = errors.New("profile not found")

// ErrAlreadyExists is returned when the destination profile name is taken.
var ErrAlreadyExists = errors.New("profile already exists")

// ErrSameName is returned when source and destination names are identical.
var ErrSameName = errors.New("source and destination names are the same")

// Store is the minimal interface required by the Renamer.
type Store interface {
	List() []string
	Get(name string) (map[string]string, error)
	Set(name string, vars map[string]string) error
	Delete(name string) error
}

// Renamer handles renaming profiles in the store.
type Renamer struct {
	store Store
}

// NewRenamer constructs a Renamer backed by the given Store.
func NewRenamer(s Store) *Renamer {
	return &Renamer{store: s}
}

// Rename copies the vars from src to dst and removes src.
// Returns an error if src does not exist, dst already exists, or src == dst.
func (r *Renamer) Rename(src, dst string) error {
	if src == dst {
		return ErrSameName
	}

	vars, err := r.store.Get(src)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, src)
	}

	existing := r.store.List()
	for _, name := range existing {
		if name == dst {
			return fmt.Errorf("%w: %s", ErrAlreadyExists, dst)
		}
	}

	if err := r.store.Set(dst, vars); err != nil {
		return fmt.Errorf("rename: create dst profile: %w", err)
	}

	if err := r.store.Delete(src); err != nil {
		// Best-effort rollback
		_ = r.store.Delete(dst)
		return fmt.Errorf("rename: remove src profile: %w", err)
	}

	return nil
}
