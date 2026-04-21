package envstats

import (
	"strings"
	"testing"
)

func TestComputeEmpty(t *testing.T) {
	s := Compute(map[string]string{})
	if s.Total != 0 {
		t.Fatalf("expected Total=0, got %d", s.Total)
	}
}

func TestComputeBasic(t *testing.T) {
	env := map[string]string{
		"HOME":  "/home/user",
		"SHELL": "/bin/zsh",
		"DEBUG": "",
	}
	s := Compute(env)

	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("Empty: want 1, got %d", s.Empty)
	}
	if s.NonEmpty != 2 {
		t.Errorf("NonEmpty: want 2, got %d", s.NonEmpty)
	}
}

func TestLongestKey(t *testing.T) {
	env := map[string]string{
		"A":              "x",
		"LONG_KEY_NAME":  "y",
		"MED":            "z",
	}
	s := Compute(env)
	if s.LongestKey != "LONG_KEY_NAME" {
		t.Errorf("LongestKey: want LONG_KEY_NAME, got %s", s.LongestKey)
	}
	if s.ShortestKey != "A" {
		t.Errorf("ShortestKey: want A, got %s", s.ShortestKey)
	}
}

func TestLongestVal(t *testing.T) {
	env := map[string]string{
		"K1": "short",
		"K2": "a much longer value here",
		"K3": "mid",
	}
	s := Compute(env)
	if s.LongestVal != "a much longer value here" {
		t.Errorf("LongestVal: want 'a much longer value here', got %q", s.LongestVal)
	}
}

func TestAvgKeyLen(t *testing.T) {
	env := map[string]string{
		"AB":   "v",
		"ABCD": "v",
	}
	s := Compute(env)
	// (2+4)/2 = 3.0
	if s.AvgKeyLen != 3.0 {
		t.Errorf("AvgKeyLen: want 3.0, got %f", s.AvgKeyLen)
	}
}

func TestFormat(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "",
	}
	s := Compute(env)
	out := Format(s)

	for _, want := range []string{"Total", "Non-empty", "Empty", "Avg key", "Longest key", "Shortest key"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output missing %q", want)
		}
	}
}

func TestSingleEntry(t *testing.T) {
	env := map[string]string{"ONLY": "value"}
	s := Compute(env)
	if s.LongestKey != "ONLY" || s.ShortestKey != "ONLY" {
		t.Errorf("single entry: longest/shortest key mismatch")
	}
	if s.Empty != 0 || s.NonEmpty != 1 {
		t.Errorf("single entry: empty/non-empty mismatch")
	}
}
