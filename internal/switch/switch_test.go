package switch_

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/shell"
)

func newTestSwitcher(t *testing.T) (*Switcher, *profile.Manager) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.json")
	m, err := profile.NewManager(path)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/bash"
	}
	e, _ := shell.NewExporter(sh)
	return NewSwitcher(m, e), m
}

func TestSwitchSuccess(t *testing.T) {
	sw, m := newTestSwitcher(t)
	if err := m.Create("dev", map[string]string{"APP_ENV": "development", "DEBUG": "true"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	result, err := sw.Switch("dev", nil)
	if err != nil {
		t.Fatalf("Switch: %v", err)
	}
	if result.NextProfile != "dev" {
		t.Errorf("expected NextProfile=dev, got %q", result.NextProfile)
	}
	if len(result.ExportLines) == 0 {
		t.Error("expected export lines, got none")
	}
}

func TestSwitchNotFound(t *testing.T) {
	sw, _ := newTestSwitcher(t)
	_, err := sw.Switch("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestSwitchDiffPopulated(t *testing.T) {
	sw, m := newTestSwitcher(t)
	if err := m.Create("prod", map[string]string{"APP_ENV": "production"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	current := map[string]string{"APP_ENV": "development", "OLD_VAR": "bye"}
	result, err := sw.Switch("prod", current)
	if err != nil {
		t.Fatalf("Switch: %v", err)
	}
	if len(result.Diff) == 0 {
		t.Error("expected non-empty diff")
	}
}

func TestPreview(t *testing.T) {
	sw, m := newTestSwitcher(t)
	if err := m.Create("staging", map[string]string{"APP_ENV": "staging", "NEW": "1"}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	current := map[string]string{"APP_ENV": "development"}
	changes, err := sw.Preview("staging", current)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if len(changes) == 0 {
		t.Error("expected changes in preview")
	}
}

func TestPreviewNotFound(t *testing.T) {
	sw, _ := newTestSwitcher(t)
	_, err := sw.Preview("ghost", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
