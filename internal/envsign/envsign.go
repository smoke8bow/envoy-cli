// Package envsign provides HMAC-based signing and verification for
// environment variable sets, allowing consumers to detect tampering.
package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrInvalidSignature is returned when a signature does not match.
var ErrInvalidSignature = errors.New("envsign: invalid signature")

// ErrEmptySecret is returned when an empty passphrase is provided.
var ErrEmptySecret = errors.New("envsign: secret must not be empty")

// Signer signs and verifies maps of environment variables.
type Signer struct {
	secret []byte
}

// New creates a Signer using the provided secret.
func New(secret string) (*Signer, error) {
	if secret == "" {
		return nil, ErrEmptySecret
	}
	return &Signer{secret: []byte(secret)}, nil
}

// Sign computes a deterministic HMAC-SHA256 signature over the given
// environment variable map. Keys are sorted before hashing so that
// insertion order does not affect the result.
func (s *Signer) Sign(vars map[string]string) string {
	h := hmac.New(sha256.New, s.secret)
	for _, k := range sortedKeys(vars) {
		fmt.Fprintf(h, "%s=%s\n", k, vars[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Verify checks that sig matches the signature computed over vars.
// Returns ErrInvalidSignature if the check fails.
func (s *Signer) Verify(vars map[string]string, sig string) error {
	expected := s.Sign(vars)
	if !hmac.Equal([]byte(expected), []byte(strings.TrimSpace(sig))) {
		return ErrInvalidSignature
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
