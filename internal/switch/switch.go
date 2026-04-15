package switch_

import (
	"fmt"

	"github.com/user/envoy-cli/internal/diff"
	"github.com/user/envoy-cli/internal/env"
	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/shell"
)

// Result holds the outcome of a profile switch operation.
type Result struct {
	PreviousProfile string
	NextProfile     string
	Diff            []diff.Change
	ExportLines     []string
}

// Switcher orchestrates switching between environment profiles.
type Switcher struct {
	manager  *profile.Manager
	exporter *shell.Exporter
}

// NewSwitcher creates a new Switcher with the given manager and exporter.
func NewSwitcher(m *profile.Manager, e *shell.Exporter) *Switcher {
	return &Switcher{manager: m, exporter: e}
}

// Switch transitions from the current environment to the named profile.
// current is a snapshot of the running environment (may be nil).
func (s *Switcher) Switch(name string, current map[string]string) (*Result, error) {
	p, err := s.manager.Get(name)
	if err != nil {
		return nil, fmt.Errorf("switch: profile %q not found: %w", name, err)
	}

	changes := diff.Compute(current, p.Vars)
	lines := s.exporter.ExportStatements(p.Vars)

	return &Result{
		PreviousProfile: env.Snapshot()["ENVOY_PROFILE"],
		NextProfile:     name,
		Diff:            changes,
		ExportLines:     lines,
	}, nil
}

// Preview returns the diff between the current environment and the named profile
// without producing export statements.
func (s *Switcher) Preview(name string, current map[string]string) ([]diff.Change, error) {
	p, err := s.manager.Get(name)
	if err != nil {
		return nil, fmt.Errorf("switch: profile %q not found: %w", name, err)
	}
	return diff.Compute(current, p.Vars), nil
}
