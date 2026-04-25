package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Entry holds the hash result for a single profile.
type Entry struct {
	Profile string
	Hash    string
}

// Hasher computes deterministic SHA-256 hashes over env var maps.
type Hasher struct{}

// New returns a new Hasher.
func New() *Hasher {
	return &Hasher{}
}

// Compute returns a deterministic hex-encoded SHA-256 hash of the given
// environment variable map. Keys are sorted before hashing so that insertion
// order does not affect the result.
func (h *Hasher) Compute(vars map[string]string) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, vars[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// Equal reports whether two env var maps produce the same hash.
func (h *Hasher) Equal(a, b map[string]string) bool {
	return h.Compute(a) == h.Compute(b)
}

// ComputeAll returns a slice of Entry values for each named profile map
// provided, sorted by profile name.
func (h *Hasher) ComputeAll(profiles map[string]map[string]string) []Entry {
	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	entries := make([]Entry, 0, len(names))
	for _, name := range names {
		entries = append(entries, Entry{
			Profile: name,
			Hash:    h.Compute(profiles[name]),
		})
	}
	return entries
}
