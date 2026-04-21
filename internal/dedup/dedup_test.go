package dedup

import (
	"testing"
)

func TestFindNoDuplicates(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	results := Find(vars)
	if len(results) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(results))
	}
}

func TestFindDuplicates(t *testing.T) {
	vars := map[string]string{"A": "same", "B": "same", "C": "other"}
	results := Find(vars)
	if len(results) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(results))
	}
	if results[0].Value != "same" {
		t.Errorf("unexpected value %q", results[0].Value)
	}
	if len(results[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(results[0].Keys))
	}
}

func TestNewInvalidStrategy(t *testing.T) {
	_, err := New("bogus")
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestApplyKeepFirst(t *testing.T) {
	vars := map[string]string{"A": "v", "B": "v", "C": "unique"}
	d, _ := New(StrategyKeepFirst)
	out := d.Apply(vars)
	// keys are sorted: A, B — keep first (A), remove B
	if _, ok := out["A"]; !ok {
		t.Error("expected A to be kept")
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if out["C"] != "unique" {
		t.Error("expected C to be unchanged")
	}
}

func TestApplyKeepLast(t *testing.T) {
	vars := map[string]string{"A": "v", "B": "v", "C": "unique"}
	d, _ := New(StrategyKeepLast)
	out := d.Apply(vars)
	// keys sorted: A, B — keep last (B), remove A
	if _, ok := out["B"]; !ok {
		t.Error("expected B to be kept")
	}
	if _, ok := out["A"]; ok {
		t.Error("expected A to be removed")
	}
}

func TestApplyRemoveAll(t *testing.T) {
	vars := map[string]string{"A": "v", "B": "v", "C": "unique"}
	d, _ := New(StrategyRemoveAll)
	out := d.Apply(vars)
	if _, ok := out["A"]; ok {
		t.Error("expected A to be removed")
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be removed")
	}
	if out["C"] != "unique" {
		t.Error("expected C to be kept")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	vars := map[string]string{"X": "dup", "Y": "dup"}
	d, _ := New(StrategyRemoveAll)
	_ = d.Apply(vars)
	if len(vars) != 2 {
		t.Error("input map was mutated")
	}
}

func TestSupported(t *testing.T) {
	s := Supported()
	if len(s) != 3 {
		t.Errorf("expected 3 strategies, got %d", len(s))
	}
}
