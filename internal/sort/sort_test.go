package sort

import (
	"testing"
)

var testVars = map[string]string{
	"ZEBRA":     "1",
	"APPLE":     "2",
	"MANGO":     "3",
	"DB_HOST":   "4",
	"A":         "5",
}

func TestIsSupportedValid(t *testing.T) {
	for _, s := range Supported() {
		if !IsSupported(s) {
			t.Errorf("expected %q to be supported", s)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if IsSupported(Strategy("bogus")) {
		t.Error("expected bogus to be unsupported")
	}
}

func TestApplyAlpha(t *testing.T) {
	keys, err := Apply(testVars, StrategyAlpha)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("not sorted alpha: %v", keys)
		}
	}
}

func TestApplyReverse(t *testing.T) {
	keys, err := Apply(testVars, StrategyReverse)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] > keys[i-1] {
			t.Errorf("not sorted reverse: %v", keys)
		}
	}
}

func TestApplyLength(t *testing.T) {
	keys, err := Apply(testVars, StrategyLength)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(keys); i++ {
		if len(keys[i]) < len(keys[i-1]) {
			t.Errorf("not sorted by length: %v", keys)
		}
	}
}

func TestApplyUnknownStrategy(t *testing.T) {
	_, err := Apply(testVars, Strategy("unknown"))
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestApplyEmptyMap(t *testing.T) {
	keys, err := Apply(map[string]string{}, StrategyAlpha)
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
