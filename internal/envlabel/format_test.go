package envlabel_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/envlabel"
)

func TestFormatEmpty(t *testing.T) {
	opts := envlabel.DefaultFormatOptions()
	got := envlabel.Format("prod", nil, opts)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestFormatBasic(t *testing.T) {
	labels := []envlabel.Label{
		{Key: "env", Value: "staging"},
		{Key: "team", Value: "backend"},
	}
	opts := envlabel.DefaultFormatOptions()
	got := envlabel.Format("prod", labels, opts)
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), got)
	}
	if lines[0] != "env=staging" {
		t.Errorf("line 0: got %q", lines[0])
	}
	if lines[1] != "team=backend" {
		t.Errorf("line 1: got %q", lines[1])
	}
}

func TestFormatShowProfile(t *testing.T) {
	labels := []envlabel.Label{
		{Key: "owner", Value: "alice"},
	}
	opts := envlabel.DefaultFormatOptions()
	opts.ShowProfile = true
	got := envlabel.Format("dev", labels, opts)
	if !strings.HasPrefix(got, "dev: ") {
		t.Errorf("expected profile prefix, got %q", got)
	}
}

func TestFormatCustomSeparator(t *testing.T) {
	labels := []envlabel.Label{
		{Key: "tier", Value: "gold"},
	}
	opts := envlabel.DefaultFormatOptions()
	opts.Separator = ":"
	got := envlabel.Format("prod", labels, opts)
	if got != "tier:gold" {
		t.Errorf("unexpected output: %q", got)
	}
}
