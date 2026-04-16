package watch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

// Watcher detects changes in a profile's env vars by comparing checksums.
type Watcher struct {
	store checksumStore
}

type checksumStore interface {
	Get(name string) (map[string]string, error)
}

// ChangeStatus represents whether a profile has changed since last snapshot.
type ChangeStatus struct {
	Profile  string
	Changed  bool
	Checksum string
}

// NewWatcher creates a Watcher backed by the given store.
func NewWatcher(s checksumStore) *Watcher {
	return &Watcher{store: s}
}

// Check computes the current checksum for profile and compares to previous.
// If previous is empty string, Changed is always false (baseline).
func (w *Watcher) Check(profile, previous string) (ChangeStatus, error) {
	vars, err := w.store.Get(profile)
	if err != nil {
		return ChangeStatus{}, err
	}
	sum := Checksum(vars)
	return ChangeStatus{
		Profile:  profile,
		Changed:  previous != "" && previous != sum,
		Checksum: sum,
	}, nil
}

// Checksum returns a deterministic SHA-256 hex digest of the env map.
func Checksum(vars map[string]string) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make([][2]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, [2]string{k, vars[k]})
	}

	b, _ := json.Marshal(ordered)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
