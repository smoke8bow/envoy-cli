package envchain_test

import (
	"errors"
	"testing"

	"github.com/user/envoy-cli/internal/envchain"
)

// fakeGetter satisfies ProfileGetter for tests.
type fakeGetter struct {
	profiles map[string]map[string]string
}

func (f *fakeGetter) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found: " + name)
	}
	return v, nil
}

func newGetter() *fakeGetter {
	return &fakeGetter{
		profiles: map[string]map[string]string{
			"base": {"HOST": "localhost", "PORT": "5432", "DEBUG": "false"},
			"prod": {"HOST": "prod.example.com", "PORT": "5432"},
			"local": {"DEBUG": "true", "PORT": "9999"},
		},
	}
}

func TestNewChainEmpty(t *testing.T) {
	_, err := envchain.NewChain(newGetter(), nil)
	if err == nil {
		t.Fatal("expected error for empty profile list")
	}
}

func TestNewChainNotFound(t *testing.T) {
	_, err := envchain.NewChain(newGetter(), []string{"missing"})
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestResolveSingleProfile(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved := c.Resolve()
	if resolved["HOST"] != "localhost" {
		t.Errorf("HOST: got %q, want %q", resolved["HOST"], "localhost")
	}
}

func TestResolveLaterProfileWins(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved := c.Resolve()
	if resolved["HOST"] != "prod.example.com" {
		t.Errorf("HOST: got %q, want %q", resolved["HOST"], "prod.example.com")
	}
	// base key not overridden
	if resolved["DEBUG"] != "false" {
		t.Errorf("DEBUG: got %q, want %q", resolved["DEBUG"], "false")
	}
}

func TestResolveThreeProfiles(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base", "prod", "local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resolved := c.Resolve()
	// local overrides prod's PORT
	if resolved["PORT"] != "9999" {
		t.Errorf("PORT: got %q, want %q", resolved["PORT"], "9999")
	}
	if resolved["DEBUG"] != "true" {
		t.Errorf("DEBUG: got %q, want %q", resolved["DEBUG"], "true")
	}
}

func TestSourceReturnsHighestPriorityProfile(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Source("HOST"); got != "prod" {
		t.Errorf("Source(HOST): got %q, want %q", got, "prod")
	}
	if got := c.Source("DEBUG"); got != "base" {
		t.Errorf("Source(DEBUG): got %q, want %q", got, "base")
	}
}

func TestSourceMissingKey(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Source("NONEXISTENT"); got != "" {
		t.Errorf("Source(NONEXISTENT): got %q, want empty string", got)
	}
}

func TestLinksCount(t *testing.T) {
	c, err := envchain.NewChain(newGetter(), []string{"base", "local"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(c.Links()); got != 2 {
		t.Errorf("Links count: got %d, want 2", got)
	}
}
