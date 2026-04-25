package envhash

import (
	"testing"
)

func newHasher() *Hasher { return New() }

func TestComputeEmptyMap(t *testing.T) {
	h := newHasher()
	hash := h.Compute(map[string]string{})
	if hash == "" {
		t.Fatal("expected non-empty hash for empty map")
	}
}

func TestComputeDeterministic(t *testing.T) {
	h := newHasher()
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if h.Compute(vars) != h.Compute(vars) {
		t.Fatal("expected identical hashes for same input")
	}
}

func TestComputeOrderIndependent(t *testing.T) {
	h := newHasher()
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"C": "3", "A": "1", "B": "2"}
	if h.Compute(a) != h.Compute(b) {
		t.Fatal("expected same hash regardless of insertion order")
	}
}

func TestComputeDifferentValues(t *testing.T) {
	h := newHasher()
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	if h.Compute(a) == h.Compute(b) {
		t.Fatal("expected different hashes for different values")
	}
}

func TestEqualTrue(t *testing.T) {
	h := newHasher()
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "1"}
	if !h.Equal(a, b) {
		t.Fatal("expected Equal to return true")
	}
}

func TestEqualFalse(t *testing.T) {
	h := newHasher()
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	if h.Equal(a, b) {
		t.Fatal("expected Equal to return false")
	}
}

func TestComputeAllSortedByName(t *testing.T) {
	h := newHasher()
	profiles := map[string]map[string]string{
		"prod":    {"ENV": "production"},
		"dev":     {"ENV": "development"},
		"staging": {"ENV": "staging"},
	}
	entries := h.ComputeAll(profiles)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	expected := []string{"dev", "prod", "staging"}
	for i, e := range entries {
		if e.Profile != expected[i] {
			t.Errorf("entry %d: expected profile %q, got %q", i, expected[i], e.Profile)
		}
		if e.Hash == "" {
			t.Errorf("entry %d: expected non-empty hash", i)
		}
	}
}

func TestComputeAllDistinctHashes(t *testing.T) {
	h := newHasher()
	profiles := map[string]map[string]string{
		"a": {"K": "v1"},
		"b": {"K": "v2"},
	}
	entries := h.ComputeAll(profiles)
	if entries[0].Hash == entries[1].Hash {
		t.Fatal("expected distinct hashes for distinct profiles")
	}
}
