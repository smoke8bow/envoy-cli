package redact

import (
	"testing"
)

func TestIsSensitiveDefault(t *testing.T) {
	r := New()
	sensitive := []string{"PASSWORD", "db_password", "API_KEY", "secret", "AUTH_TOKEN", "private_key"}
	for _, k := range sensitive {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsNotSensitive(t *testing.T) {
	r := New()
	safe := []string{"HOST", "PORT", "DEBUG", "APP_NAME"}
	for _, k := range safe {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to not be sensitive", k)
		}
	}
}

func TestApplyRedacts(t *testing.T) {
	r := New()
	vars := map[string]string{
		"HOST":     "localhost",
		"PASSWORD": "s3cr3t",
		"API_KEY":  "abc123",
	}
	out := r.Apply(vars)
	if out["HOST"] != "localhost" {
		t.Errorf("HOST should not be redacted")
	}
	if out["PASSWORD"] != "[REDACTED]" {
		t.Errorf("PASSWORD should be redacted, got %q", out["PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", out["API_KEY"])
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	r := New()
	vars := map[string]string{"PASSWORD": "original"}
	_ = r.Apply(vars)
	if vars["PASSWORD"] != "original" {
		t.Error("Apply must not mutate input map")
	}
}

func TestRedactString(t *testing.T) {
	r := New()
	vars := map[string]string{"API_KEY": "supersecret"}
	s := "Authorization: Bearer supersecret"
	out := r.RedactString(s, vars)
	if out != "Authorization: Bearer [REDACTED]" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestNewWithPatternsInvalid(t *testing.T) {
	_, err := NewWithPatterns([]string{"["})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestNewWithPatternsCustom(t *testing.T) {
	r, err := NewWithPatterns([]string{`(?i)credential`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsSensitive("DB_CREDENTIAL") {
		t.Error("expected DB_CREDENTIAL to be sensitive")
	}
	if r.IsSensitive("PASSWORD") {
		t.Error("custom redactor should not include default patterns")
	}
}
