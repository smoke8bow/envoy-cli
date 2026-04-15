package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/clone"
	"github.com/envoy-cli/envoy/internal/profile"
	"github.com/envoy-cli/envoy/internal/store"
)

func newCloner(t *testing.T) *clone.Cloner {
	t.Helper()
	dir, err := os.MkdirTemp("", "clone-test-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	s, err := store.Load(filepath.Join(dir, "store.json"))
	if err != nil {
		t.Fatalf("store.Load: %v", err)
	}
	m := profile.NewManager(s)
	return clone.NewCloner(m)
}

func TestCloneSuccess(t *testing.T) {
	c := newCloner(t)

	// seed source
	if err := c.Manager().Create("prod", map[string]string{"HOST": "prod.example.com", "PORT": "443"}); err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := c.Clone("prod", "staging"); err != nil {
		t.Fatalf("Clone: %v", err)
	}

	p, err := c.Manager().Get("staging")
	if err != nil {
		t.Fatalf("get staging: %v", err)
	}
	if p.Vars["HOST"] != "prod.example.com" {
		t.Errorf("HOST = %q, want prod.example.com", p.Vars["HOST"])
	}
}

func TestCloneSourceNotFound(t *testing.T) {
	c := newCloner(t)
	if err := c.Clone("ghost", "copy"); err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCloneDuplicateDst(t *testing.T) {
	c := newCloner(t)
	c.Manager().Create("a", map[string]string{"K": "v"})
	c.Manager().Create("b", map[string]string{"K": "v"})
	if err := c.Clone("a", "b"); err == nil {
		t.Fatal("expected error when dst already exists")
	}
}

func TestCloneInvalidName(t *testing.T) {
	c := newCloner(t)
	c.Manager().Create("src", map[string]string{"X": "1"})
	if err := c.Clone("src", "bad name!"); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCloneWithOverrides(t *testing.T) {
	c := newCloner(t)
	c.Manager().Create("base", map[string]string{"HOST": "base.local", "PORT": "80"})

	err := c.CloneWithOverrides("base", "dev", map[string]string{"HOST": "dev.local"})
	if err != nil {
		t.Fatalf("CloneWithOverrides: %v", err)
	}

	p, _ := c.Manager().Get("dev")
	if p.Vars["HOST"] != "dev.local" {
		t.Errorf("HOST = %q, want dev.local", p.Vars["HOST"])
	}
	if p.Vars["PORT"] != "80" {
		t.Errorf("PORT = %q, want 80", p.Vars["PORT"])
	}
}

func TestCloneIsolation(t *testing.T) {
	c := newCloner(t)
	c.Manager().Create("original", map[string]string{"KEY": "original"})
	c.Clone("original", "copy")

	c.Manager().Update("copy", map[string]string{"KEY": "modified"})

	orig, _ := c.Manager().Get("original")
	if orig.Vars["KEY"] != "original" {
		t.Errorf("original mutated: KEY = %q", orig.Vars["KEY"])
	}
}
