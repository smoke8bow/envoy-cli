package merge

import (
	"testing"
)

func TestMergeNoConflict(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}

	r, err := Merge(dst, src, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Vars["A"] != "1" || r.Vars["B"] != "2" {
		t.Errorf("unexpected vars: %v", r.Vars)
	}
	if len(r.Added) != 1 || r.Added[0] != "B" {
		t.Errorf("expected B in Added, got %v", r.Added)
	}
}

func TestMergeStrategyOurs(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}

	r, err := Merge(dst, src, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Vars["A"] != "original" {
		t.Errorf("expected original value, got %q", r.Vars["A"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "A" {
		t.Errorf("expected A in Skipped, got %v", r.Skipped)
	}
}

func TestMergeStrategyTheirs(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}

	r, err := Merge(dst, src, StrategyTheirs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Vars["A"] != "new" {
		t.Errorf("expected new value, got %q", r.Vars["A"])
	}
	if len(r.Overwrite) != 1 || r.Overwrite[0] != "A" {
		t.Errorf("expected A in Overwrite, got %v", r.Overwrite)
	}
}

func TestMergeStrategyError(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "2"}

	_, err := Merge(dst, src, StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMergeDoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}

	_, err := Merge(dst, src, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["B"]; ok {
		t.Error("dst was mutated")
	}
}

func TestMergeEmptySrc(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{}

	r, err := Merge(dst, src, StrategyOurs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Vars) != 1 || r.Vars["A"] != "1" {
		t.Errorf("unexpected vars: %v", r.Vars)
	}
	if len(r.Added) != 0 || len(r.Skipped) != 0 || len(r.Overwrite) != 0 {
		t.Errorf("expected empty slices for empty src")
	}
}
