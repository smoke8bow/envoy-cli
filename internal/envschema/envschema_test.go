package envschema_test

import (
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/envschema"
)

func newSchema(fields ...envschema.Field) *envschema.Schema {
	return &envschema.Schema{Fields: fields}
}

func TestNoViolationsAllPresent(t *testing.T) {
	s := newSchema(
		envschema.Field{Key: "HOST", Type: envschema.TypeString, Required: true},
		envschema.Field{Key: "PORT", Type: envschema.TypeInt, Required: true},
	)
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	if v := s.Validate(env); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestRequiredMissingKey(t *testing.T) {
	s := newSchema(envschema.Field{Key: "SECRET", Type: envschema.TypeString, Required: true})
	env := map[string]string{}
	v := s.Validate(env)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "SECRET" {
		t.Errorf("unexpected key %q", v[0].Key)
	}
}

func TestRequiredEmptyValue(t *testing.T) {
	s := newSchema(envschema.Field{Key: "TOKEN", Type: envschema.TypeString, Required: true})
	env := map[string]string{"TOKEN": "   "}
	v := s.Validate(env)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestOptionalMissingKeyNoViolation(t *testing.T) {
	s := newSchema(envschema.Field{Key: "OPTIONAL", Type: envschema.TypeString, Required: false})
	env := map[string]string{}
	if v := s.Validate(env); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestTypeIntValid(t *testing.T) {
	s := newSchema(envschema.Field{Key: "COUNT", Type: envschema.TypeInt, Required: true})
	env := map[string]string{"COUNT": "42"}
	if v := s.Validate(env); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestTypeIntInvalid(t *testing.T) {
	s := newSchema(envschema.Field{Key: "COUNT", Type: envschema.TypeInt, Required: true})
	env := map[string]string{"COUNT": "notanint"}
	v := s.Validate(env)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestTypeBoolValid(t *testing.T) {
	s := newSchema(envschema.Field{Key: "ENABLED", Type: envschema.TypeBool, Required: true})
	for _, val := range []string{"true", "false", "1", "0"} {
		env := map[string]string{"ENABLED": val}
		if v := s.Validate(env); len(v) != 0 {
			t.Errorf("expected no violations for %q, got %v", val, v)
		}
	}
}

func TestTypeBoolInvalid(t *testing.T) {
	s := newSchema(envschema.Field{Key: "ENABLED", Type: envschema.TypeBool, Required: true})
	env := map[string]string{"ENABLED": "yes"}
	v := s.Validate(env)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestTypeFloatValid(t *testing.T) {
	s := newSchema(envschema.Field{Key: "RATIO", Type: envschema.TypeFloat, Required: true})
	env := map[string]string{"RATIO": "3.14"}
	if v := s.Validate(env); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestViolationErrorString(t *testing.T) {
	v := envschema.Violation{Key: "FOO", Message: "required key is missing or empty"}
	if v.Error() == "" {
		t.Error("expected non-empty error string")
	}
}
