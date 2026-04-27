// Package envschema provides schema validation for environment variable profiles.
// It allows defining expected keys with type constraints and optional/required flags.
package envschema

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FieldType represents the expected type of an environment variable value.
type FieldType string

const (
	TypeString FieldType = "string"
	TypeInt    FieldType = "int"
	TypeBool   FieldType = "bool"
	TypeFloat  FieldType = "float"
)

// Field describes a single expected environment variable.
type Field struct {
	Key      string
	Type     FieldType
	Required bool
}

// Schema holds a collection of field definitions.
type Schema struct {
	Fields []Field
}

// Violation describes a single schema validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks the given env map against the schema.
// It returns a list of violations; an empty slice means the map is valid.
func (s *Schema) Validate(env map[string]string) []Violation {
	var violations []Violation

	for _, f := range s.Fields {
		val, ok := env[f.Key]
		if !ok || strings.TrimSpace(val) == "" {
			if f.Required {
				violations = append(violations, Violation{Key: f.Key, Message: "required key is missing or empty"})
			}
			continue
		}
		if err := checkType(val, f.Type); err != nil {
			violations = append(violations, Violation{Key: f.Key, Message: err.Error()})
		}
	}

	return violations
}

func checkType(val string, t FieldType) error {
	switch t {
	case TypeInt:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return fmt.Errorf("expected int, got %q", val)
		}
	case TypeBool:
		if _, err := strconv.ParseBool(val); err != nil {
			return fmt.Errorf("expected bool, got %q", val)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Errorf("expected float, got %q", val)
		}
	case TypeString:
		// any value is valid
	default:
		return errors.New("unknown field type")
	}
	return nil
}
