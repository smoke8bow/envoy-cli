package envmerge

import (
	"testing"
)

func TestIsSupportedValid(t *testing.T) {
	for _, s := range Supported() {
		if !IsSupported(s) {
			t.Errorf("expected %q to be supported", s)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if IsSupported("unknown") {
		t.Error("expected 'unknown' to be unsupported")
	}
}

func TestMergeNoConflict(t *testing.T) {
	sources := []Source{
		{Name: "base", Vars: map[string]string{"A": "1", "B": "2"}},
		{Name: "extra", Vars: map[string]string{"C": "3"}},
	}
	res, err := Merge(sources, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["A"] != "1" || res.Vars["B"] != "2" || res.Vars["C"] != "3" {
		t.Errorf("unexpected vars: %v", res.Vars)
	}
	if res.Origin["C"] != "extra" {
		t.Errorf("expected origin 'extra', got %q", res.Origin["C"])
	}
}

func TestMergeStrategyFirst(t *testing.T) {
	sources := []Source{
		{Name: "base", Vars: map[string]string{"KEY": "original"}},
		{Name: "override", Vars: map[string]string{"KEY": "new"}},
	}
	res, err := Merge(sources, StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "original" {
		t.Errorf("expected 'original', got %q", res.Vars["KEY"])
	}
	if res.Origin["KEY"] != "base" {
		t.Errorf("expected origin 'base', got %q", res.Origin["KEY"])
	}
}

func TestMergeStrategyLast(t *testing.T) {
	sources := []Source{
		{Name: "base", Vars: map[string]string{"KEY": "original"}},
		{Name: "override", Vars: map[string]string{"KEY": "new"}},
	}
	res, err := Merge(sources, StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "new" {
		t.Errorf("expected 'new', got %q", res.Vars["KEY"])
	}
	if res.Origin["KEY"] != "override" {
		t.Errorf("expected origin 'override', got %q", res.Origin["KEY"])
	}
}

func TestMergeStrategyStrict(t *testing.T) {
	sources := []Source{
		{Name: "a", Vars: map[string]string{"KEY": "1"}},
		{Name: "b", Vars: map[string]string{"KEY": "2"}},
	}
	_, err := Merge(sources, StrategyStrict)
	if err == nil {
		t.Error("expected conflict error, got nil")
	}
}

func TestMergeStrictNoConflict(t *testing.T) {
	sources := []Source{
		{Name: "a", Vars: map[string]string{"X": "1"}},
		{Name: "b", Vars: map[string]string{"Y": "2"}},
	}
	res, err := Merge(sources, StrategyStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(res.Vars))
	}
}

func TestMergeUnsupportedStrategy(t *testing.T) {
	_, err := Merge(nil, Strategy("bogus"))
	if err == nil {
		t.Error("expected error for unsupported strategy")
	}
}

func TestMergeEmptySources(t *testing.T) {
	res, err := Merge([]Source{}, StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty result, got %v", res.Vars)
	}
}
