package envmap_test

import (
	"testing"

	"envoy-cli/internal/envmap"
)

func TestFromSliceBasic(t *testing.T) {
	m := envmap.FromSlice([]string{"FOO=bar", "BAZ=qux"})
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"])
	}
}

func TestFromSliceNoEquals(t *testing.T) {
	m := envmap.FromSlice([]string{"NOVALUE"})
	v, ok := m["NOVALUE"]
	if !ok {
		t.Fatal("expected key NOVALUE to exist")
	}
	if v != "" {
		t.Errorf("expected empty value, got %q", v)
	}
}

func TestFromSliceValueContainsEquals(t *testing.T) {
	m := envmap.FromSlice([]string{"URL=http://x.com?a=1&b=2"})
	if m["URL"] != "http://x.com?a=1&b=2" {
		t.Errorf("unexpected value: %q", m["URL"])
	}
}

func TestToSliceSorted(t *testing.T) {
	m := map[string]string{"Z": "last", "A": "first", "M": "mid"}
	s := envmap.ToSlice(m)
	expected := []string{"A=first", "M=mid", "Z=last"}
	for i, v := range expected {
		if s[i] != v {
			t.Errorf("index %d: expected %q got %q", i, v, s[i])
		}
	}
}

func TestKeysReturnsAllSorted(t *testing.T) {
	m := map[string]string{"C": "3", "A": "1", "B": "2"}
	keys := envmap.Keys(m)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestCloneIsIndependent(t *testing.T) {
	orig := map[string]string{"FOO": "original"}
	cloned := envmap.Clone(orig)
	cloned["FOO"] = "changed"
	if orig["FOO"] != "original" {
		t.Error("Clone mutated original map")
	}
}

func TestMergeOverwrites(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}
	envmap.Merge(dst, src)
	if dst["B"] != "99" {
		t.Errorf("expected B=99, got %q", dst["B"])
	}
	if dst["C"] != "3" {
		t.Errorf("expected C=3, got %q", dst["C"])
	}
	if dst["A"] != "1" {
		t.Errorf("expected A=1 unchanged, got %q", dst["A"])
	}
}

func TestMergeDoesNotMutateSrc(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"X": "10"}
	envmap.Merge(dst, src)
	dst["X"] = "changed"
	if src["X"] != "10" {
		t.Error("Merge mutated src map")
	}
}
