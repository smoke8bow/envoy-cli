package diff2_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/diff2"
)

func TestFormatNoChanges(t *testing.T) {
	changes := []diff2.Change{}
	out := diff2.Format(changes, diff2.DefaultFormatOptions())
	if out != "" {
		t.Errorf("expected empty output for no changes, got %q", out)
	}
}

func TestFormatAdded(t *testing.T) {
	changes := []diff2.Change{
		{Key: "NEW_KEY", Right: "value", Kind: diff2.Added},
	}
	out := diff2.Format(changes, diff2.DefaultFormatOptions())
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected output to contain key, got %q", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected output to contain '+' marker for added key, got %q", out)
	}
}

func TestFormatRemoved(t *testing.T) {
	changes := []diff2.Change{
		{Key: "OLD_KEY", Left: "value", Kind: diff2.Removed},
	}
	out := diff2.Format(changes, diff2.DefaultFormatOptions())
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected output to contain key, got %q", out)
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected output to contain '-' marker for removed key, got %q", out)
	}
}

func TestFormatChanged(t *testing.T) {
	changes := []diff2.Change{
		{Key: "DB_HOST", Left: "localhost", Right: "prod.db.internal", Kind: diff2.Changed},
	}
	out := diff2.Format(changes, diff2.DefaultFormatOptions())
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected output to contain key")
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected output to contain old value")
	}
	if !strings.Contains(out, "prod.db.internal") {
		t.Errorf("expected output to contain new value")
	}
}

func TestSummaryCountsCorrect(t *testing.T) {
	changes := []diff2.Change{
		{Key: "A", Right: "1", Kind: diff2.Added},
		{Key: "B", Right: "2", Kind: diff2.Added},
		{Key: "C", Left: "old", Kind: diff2.Removed},
		{Key: "D", Left: "x", Right: "y", Kind: diff2.Changed},
		{Key: "E", Left: "same", Right: "same", Kind: diff2.Unchanged},
	}
	summary := diff2.Summary(changes)
	if summary.Added != 2 {
		t.Errorf("expected 2 added, got %d", summary.Added)
	}
	if summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", summary.Removed)
	}
	if summary.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", summary.Changed)
	}
	if summary.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", summary.Unchanged)
	}
}

func TestFormatHideUnchanged(t *testing.T) {
	changes := []diff2.Change{
		{Key: "STABLE", Left: "v", Right: "v", Kind: diff2.Unchanged},
		{Key: "CHANGED", Left: "a", Right: "b", Kind: diff2.Changed},
	}
	opts := diff2.DefaultFormatOptions()
	opts.ShowUnchanged = false
	out := diff2.Format(changes, opts)
	if strings.Contains(out, "STABLE") {
		t.Errorf("expected unchanged key to be hidden, got %q", out)
	}
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected changed key to appear in output")
	}
}

func TestFormatShowUnchanged(t *testing.T) {
	changes := []diff2.Change{
		{Key: "STABLE", Left: "v", Right: "v", Kind: diff2.Unchanged},
	}
	opts := diff2.DefaultFormatOptions()
	opts.ShowUnchanged = true
	out := diff2.Format(changes, opts)
	if !strings.Contains(out, "STABLE") {
		t.Errorf("expected unchanged key to appear when ShowUnchanged=true, got %q", out)
	}
}
