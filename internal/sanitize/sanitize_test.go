package sanitize_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/sanitize"
)

func newSanitizer(opts sanitize.Options) *sanitize.Sanitizer {
	return sanitize.New(opts)
}

func TestApplyTrimSpace(t *testing.T) {
	s := newSanitizer(sanitize.Options{TrimSpace: true})
	result := s.Apply(map[string]string{
		"KEY": "  hello  ",
		"OTHER": "\tworld\t",
	})
	if result["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", result["KEY"])
	}
	if result["OTHER"] != "world" {
		t.Errorf("expected 'world', got %q", result["OTHER"])
	}
}

func TestApplyStripControl(t *testing.T) {
	s := newSanitizer(sanitize.Options{StripControl: true})
	input := "hello\x00world\x01"
	result := s.Value(input)
	if strings.ContainsAny(result, "\x00\x01") {
		t.Errorf("control characters not stripped: %q", result)
	}
	if result != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", result)
	}
}

func TestApplyTabPreservedWhenStripControl(t *testing.T) {
	s := newSanitizer(sanitize.Options{StripControl: true})
	result := s.Value("col1\tcol2")
	if result != "col1\tcol2" {
		t.Errorf("tab should be preserved, got %q", result)
	}
}

func TestApplyMaxValueLen(t *testing.T) {
	s := newSanitizer(sanitize.Options{MaxValueLen: 5})
	result := s.Value("abcdefgh")
	if result != "abcde" {
		t.Errorf("expected 'abcde', got %q", result)
	}
}

func TestApplyMaxValueLenNoTruncation(t *testing.T) {
	s := newSanitizer(sanitize.Options{MaxValueLen: 10})
	result := s.Value("short")
	if result != "short" {
		t.Errorf("expected 'short', got %q", result)
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	s := newSanitizer(sanitize.DefaultOptions())
	input := map[string]string{"K": "  v  "}
	_ = s.Apply(input)
	if input["K"] != "  v  " {
		t.Error("input map was mutated")
	}
}

func TestDefaultOptionsApplied(t *testing.T) {
	s := sanitize.New(sanitize.DefaultOptions())
	result := s.Apply(map[string]string{
		"A": "  \x00hello\x00  ",
	})
	if result["A"] != "hello" {
		t.Errorf("expected 'hello', got %q", result["A"])
	}
}
