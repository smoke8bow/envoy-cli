package mask

import (
	"testing"
)

func TestIsSensitiveDefault(t *testing.T) {
	m := NewMasker(nil)
	cases := []struct {
		key  string
		want bool
	}{
		{"DB_PASSWORD", true},
		{"API_TOKEN", true},
		{"SECRET_KEY", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, c := range cases {
		if got := m.IsSensitive(c.key); got != c.want {
			t.Errorf("IsSensitive(%q) = %v, want %v", c.key, got, c.want)
		}
	}
}

func TestMaskValue(t *testing.T) {
	m := NewMasker(nil)
	if got := m.MaskValue("API_TOKEN", "abc123"); got != "***" {
		t.Errorf("expected *** got %q", got)
	}
	if got := m.MaskValue("APP_NAME", "myapp"); got != "myapp" {
		t.Errorf("expected myapp got %q", got)
	}
}

func TestApply(t *testing.T) {
	m := NewMasker(nil)
	input := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_TOKEN":   "tok_xyz",
	}
	out := m.Apply(input)
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be masked")
	}
	if out["PORT"] != "8080" {
		t.Errorf("PORT should not be masked")
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if out["API_TOKEN"] != "***" {
		t.Errorf("API_TOKEN should be masked")
	}
}

func TestCustomPatterns(t *testing.T) {
	m := NewMasker([]string{"CUSTOM"})
	if !m.IsSensitive("MY_CUSTOM_VAR") {
		t.Error("expected MY_CUSTOM_VAR to be sensitive")
	}
	if m.IsSensitive("API_TOKEN") {
		t.Error("API_TOKEN should not be sensitive with custom patterns")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	m := NewMasker(nil)
	input := map[string]string{"DB_PASSWORD": "secret"}
	m.Apply(input)
	if input["DB_PASSWORD"] != "secret" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestMaskValueEmptyString(t *testing.T) {
	// Masking an empty sensitive value should still return the mask token,
	// not an empty string, so callers cannot infer the value is unset.
	m := NewMasker(nil)
	if got := m.MaskValue("DB_PASSWORD", ""); got != "***" {
		t.Errorf("expected *** for empty sensitive value, got %q", got)
	}
	// An empty non-sensitive value should pass through unchanged.
	if got := m.MaskValue("APP_NAME", ""); got != "" {
		t.Errorf("expected empty string for non-sensitive key, got %q", got)
	}
}
