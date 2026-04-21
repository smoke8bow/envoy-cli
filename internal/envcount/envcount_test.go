package envcount_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/envcount"
)

func newCounter() *envcount.Counter { return envcount.New() }

func TestComputeEmpty(t *testing.T) {
	c := newCounter()
	s := c.Compute("dev", map[string]string{})
	if s.Total != 0 || s.Empty != 0 || s.NonEmpty != 0 {
		t.Fatalf("expected all zeros, got %+v", s)
	}
	if s.Name != "dev" {
		t.Fatalf("unexpected name %q", s.Name)
	}
}

func TestComputeCounts(t *testing.T) {
	c := newCounter()
	vars := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
		"KEY_C": "another",
	}
	s := c.Compute("prod", vars)
	if s.Total != 3 {
		t.Fatalf("expected total 3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Fatalf("expected empty 1, got %d", s.Empty)
	}
	if s.NonEmpty != 2 {
		t.Fatalf("expected non-empty 2, got %d", s.NonEmpty)
	}
}

func TestComputeAllSorted(t *testing.T) {
	c := newCounter()
	profiles := map[string]map[string]string{
		"zebra": {"A": "1"},
		"alpha": {"B": "2", "C": ""},
	}
	results := c.ComputeAll(profiles)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Name != "alpha" {
		t.Fatalf("expected alpha first, got %q", results[0].Name)
	}
	if results[1].Name != "zebra" {
		t.Fatalf("expected zebra second, got %q", results[1].Name)
	}
}

func TestTotals(t *testing.T) {
	stats := []envcount.ProfileStats{
		{Name: "a", Total: 3, Empty: 1, NonEmpty: 2},
		{Name: "b", Total: 5, Empty: 0, NonEmpty: 5},
	}
	tot := envcount.Totals(stats)
	if tot.Total != 8 {
		t.Fatalf("expected total 8, got %d", tot.Total)
	}
	if tot.Empty != 1 {
		t.Fatalf("expected empty 1, got %d", tot.Empty)
	}
	if tot.NonEmpty != 7 {
		t.Fatalf("expected non-empty 7, got %d", tot.NonEmpty)
	}
}

func TestProfileStatsString(t *testing.T) {
	s := envcount.ProfileStats{Name: "dev", Total: 4, NonEmpty: 3, Empty: 1}
	got := s.String()
	expected := "dev: total=4 non-empty=3 empty=1"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
