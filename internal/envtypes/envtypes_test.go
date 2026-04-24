package envtypes

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
	if IsSupported("binary") {
		t.Error("expected 'binary' to be unsupported")
	}
}

func TestInferBool(t *testing.T) {
	for _, v := range []string{"true", "false", "True", "FALSE"} {
		if got := Infer(v); got != TypeBool {
			t.Errorf("Infer(%q) = %q, want %q", v, got, TypeBool)
		}
	}
}

func TestInferInt(t *testing.T) {
	for _, v := range []string{"0", "42", "-7", "1000000"} {
		if got := Infer(v); got != TypeInt {
			t.Errorf("Infer(%q) = %q, want %q", v, got, TypeInt)
		}
	}
}

func TestInferFloat(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "1e10"} {
		if got := Infer(v); got != TypeFloat {
			t.Errorf("Infer(%q) = %q, want %q", v, got, TypeFloat)
		}
	}
}

func TestInferURL(t *testing.T) {
	for _, v := range []string{"http://example.com", "https://api.example.com/v1"} {
		if got := Infer(v); got != TypeURL {
			t.Errorf("Infer(%q) = %q, want %q", v, got, TypeURL)
		}
	}
}

func TestInferString(t *testing.T) {
	for _, v := range []string{"hello", "some-value", ""} {
		if got := Infer(v); got != TypeString {
			t.Errorf("Infer(%q) = %q, want %q", v, got, TypeString)
		}
	}
}

func TestInferAll(t *testing.T) {
	vars := map[string]string{
		"PORT":    "8080",
		"DEBUG":   "true",
		"API_URL": "https://api.example.com",
		"NAME":    "envoy",
	}
	result := InferAll(vars)
	expected := map[string]Type{
		"PORT":    TypeInt,
		"DEBUG":   TypeBool,
		"API_URL": TypeURL,
		"NAME":    TypeString,
	}
	for k, want := range expected {
		if got := result[k]; got != want {
			t.Errorf("InferAll[%s] = %q, want %q", k, got, want)
		}
	}
}

func TestValidateBoolOK(t *testing.T) {
	if err := Validate("DEBUG", "true", TypeBool); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBoolFail(t *testing.T) {
	if err := Validate("DEBUG", "yes", TypeBool); err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestValidateIntFail(t *testing.T) {
	if err := Validate("PORT", "abc", TypeInt); err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestValidateURLFail(t *testing.T) {
	if err := Validate("URL", "ftp://bad.com", TypeURL); err == nil {
		t.Error("expected error for non-http URL")
	}
}

func TestValidateStringAlwaysOK(t *testing.T) {
	if err := Validate("KEY", "anything goes", TypeString); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
