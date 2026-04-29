// Package envfence restricts which environment variable keys are allowed
// in a profile by enforcing an allowlist or denylist fence.
package envfence

import (
	"errors"
	"fmt"
	"strings"
)

// Mode controls whether the fence is an allowlist or denylist.
type Mode string

const (
	ModeAllow Mode = "allow"
	ModeDeny  Mode = "deny"
)

// Violation describes a key that violated the fence policy.
type Violation struct {
	Key    string
	Reason string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Reason)
}

// Fence holds the policy configuration.
type Fence struct {
	mode Mode
	keys map[string]struct{}
}

// New creates a Fence with the given mode and set of keys.
// mode must be ModeAllow or ModeDeny.
func New(mode Mode, keys []string) (*Fence, error) {
	if mode != ModeAllow && mode != ModeDeny {
		return nil, fmt.Errorf("unsupported mode %q: must be \"allow\" or \"deny\"", mode)
	}
	if len(keys) == 0 {
		return nil, errors.New("fence requires at least one key")
	}
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[strings.TrimSpace(k)] = struct{}{}
	}
	return &Fence{mode: mode, keys: km}, nil
}

// Check validates vars against the fence policy and returns any violations.
func (f *Fence) Check(vars map[string]string) []Violation {
	var violations []Violation
	for k := range vars {
		_, listed := f.keys[k]
		switch f.mode {
		case ModeAllow:
			if !listed {
				violations = append(violations, Violation{Key: k, Reason: "not in allowlist"})
			}
		case ModeDeny:
			if listed {
				violations = append(violations, Violation{Key: k, Reason: "present in denylist"})
			}
		}
	}
	return violations
}

// Filter returns a copy of vars with disallowed keys removed (allowlist mode)
// or denied keys removed (denylist mode).
func (f *Fence) Filter(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		_, listed := f.keys[k]
		switch f.mode {
		case ModeAllow:
			if listed {
				out[k] = v
			}
		case ModeDeny:
			if !listed {
				out[k] = v
			}
		}
	}
	return out
}
