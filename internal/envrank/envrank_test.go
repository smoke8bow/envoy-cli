package envrank

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

func TestRankByKeyLen(t *testing.T) {
	vars := map[string]string{
		"AB":     "x",
		"ABCDEF": "y",
		"ABC":    "z",
	}
	entries, err := Rank(vars, StrategyKeyLen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "ABCDEF" {
		t.Errorf("expected first key ABCDEF, got %s", entries[0].Key)
	}
	if entries[0].Rank != 6 {
		t.Errorf("expected rank 6, got %d", entries[0].Rank)
	}
}

func TestRankByValueLen(t *testing.T) {
	vars := map[string]string{
		"A": "hello world",
		"B": "hi",
		"C": "hey there!",
	}
	entries, err := Rank(vars, StrategyValueLen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "A" {
		t.Errorf("expected first key A (longest value), got %s", entries[0].Key)
	}
}

func TestRankByAlpha(t *testing.T) {
	vars := map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
	}
	entries, err := Rank(vars, StrategyAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// alphabeticScore returns -int(key[0]), so higher first byte => higher rank
	// 'Z'=90 > 'M'=77 > 'A'=65, so ZEBRA first
	if entries[0].Key != "ZEBRA" {
		t.Errorf("expected ZEBRA first, got %s", entries[0].Key)
	}
	if entries[2].Key != "APPLE" {
		t.Errorf("expected APPLE last, got %s", entries[2].Key)
	}
}

func TestRankUnsupportedStrategy(t *testing.T) {
	_, err := Rank(map[string]string{"K": "V"}, "bogus")
	if err == nil {
		t.Error("expected error for unsupported strategy")
	}
}

func TestRankEmptyMap(t *testing.T) {
	entries, err := Rank(map[string]string{}, StrategyKeyLen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRankTieBreakByKey(t *testing.T) {
	vars := map[string]string{
		"BB": "x",
		"AA": "y",
		"CC": "z",
	}
	entries, err := Rank(vars, StrategyKeyLen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All keys same length; tie-break alphabetically
	if entries[0].Key != "AA" {
		t.Errorf("expected AA first on tie-break, got %s", entries[0].Key)
	}
}
