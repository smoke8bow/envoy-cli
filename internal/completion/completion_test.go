package completion_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/your-org/envoy-cli/internal/completion"
)

// mockLister implements ProfileLister for testing.
type mockLister struct {
	names []string
	err   error
}

func (m *mockLister) List() ([]string, error) {
	return m.names, m.err
}

func TestProfileNamesJoined(t *testing.T) {
	lister := &mockLister{names: []string{"dev", "staging", "prod"}}
	g := completion.NewGenerator(completion.Bash, lister)

	out, err := g.ProfileNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "dev\nstaging\nprod" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestProfileNamesEmpty(t *testing.T) {
	lister := &mockLister{names: []string{}}
	g := completion.NewGenerator(completion.Zsh, lister)

	out, err := g.ProfileNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestProfileNamesError(t *testing.T) {
	lister := &mockLister{err: errors.New("store failure")}
	g := completion.NewGenerator(completion.Bash, lister)

	_, err := g.ProfileNames()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "store failure") {
		t.Errorf("error should mention underlying cause, got: %v", err)
	}
}

func TestScriptBash(t *testing.T) {
	g := completion.NewGenerator(completion.Bash, &mockLister{})
	script, err := g.Script("envoy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(script, "complete -F") {
		t.Errorf("bash script missing 'complete -F': %s", script)
	}
	if !strings.Contains(script, "envoy") {
		t.Errorf("bash script missing program name: %s", script)
	}
}

func TestScriptZsh(t *testing.T) {
	g := completion.NewGenerator(completion.Zsh, &mockLister{})
	script, err := g.Script("envoy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(script, "#compdef") {
		t.Errorf("zsh script missing '#compdef': %s", script)
	}
}

func TestScriptFish(t *testing.T) {
	g := completion.NewGenerator(completion.Fish, &mockLister{})
	script, err := g.Script("envoy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(script, "complete -c") {
		t.Errorf("fish script missing 'complete -c': %s", script)
	}
}

func TestScriptUnsupportedShell(t *testing.T) {
	g := completion.NewGenerator(completion.Shell("powershell"), &mockLister{})
	_, err := g.Script("envoy")
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("error should mention unsupported shell, got: %v", err)
	}
}
