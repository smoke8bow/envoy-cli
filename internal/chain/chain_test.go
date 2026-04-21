package chain_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/chain"
)

func TestComposeEmpty(t *testing.T) {
	c, err := chain.NewComposer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := c.Compose()
	if len(r.Vars) != 0 {
		t.Fatalf("expected empty vars, got %v", r.Vars)
	}
}

func TestComposeSingleLayer(t *testing.T) {
	entries := []chain.Entry{
		{Name: "base", Vars: map[string]string{"A": "1", "B": "2"}},
	}
	c, _ := chain.NewComposer(entries)
	r := c.Compose()
	if r.Vars["A"] != "1" || r.Vars["B"] != "2" {
		t.Fatalf("unexpected vars: %v", r.Vars)
	}
	if r.Source["A"] != "base" {
		t.Fatalf("expected source 'base', got %q", r.Source["A"])
	}
}

func TestComposeLaterOverrides(t *testing.T) {
	entries := []chain.Entry{
		{Name: "base", Vars: map[string]string{"A": "base-val", "B": "b"}},
		{Name: "override", Vars: map[string]string{"A": "new-val"}},
	}
	c, _ := chain.NewComposer(entries)
	r := c.Compose()
	if r.Vars["A"] != "new-val" {
		t.Fatalf("expected 'new-val', got %q", r.Vars["A"])
	}
	if r.Source["A"] != "override" {
		t.Fatalf("expected source 'override', got %q", r.Source["A"])
	}
	if r.Vars["B"] != "b" {
		t.Fatalf("expected B='b', got %q", r.Vars["B"])
	}
}

func TestComposeLayers(t *testing.T) {
	entries := []chain.Entry{
		{Name: "a", Vars: map[string]string{}},
		{Name: "b", Vars: map[string]string{}},
	}
	c, _ := chain.NewComposer(entries)
	layers := c.Layers()
	if len(layers) != 2 || layers[0] != "a" || layers[1] != "b" {
		t.Fatalf("unexpected layers: %v", layers)
	}
}

func TestNewComposerEmptyNameError(t *testing.T) {
	_, err := chain.NewComposer([]chain.Entry{
		{Name: "", Vars: map[string]string{}},
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewComposerNilVarsError(t *testing.T) {
	_, err := chain.NewComposer([]chain.Entry{
		{Name: "x", Vars: nil},
	})
	if err == nil {
		t.Fatal("expected error for nil vars")
	}
}
