package envtypes

import (
	"fmt"
	"strconv"
	"strings"
)

// Type represents the inferred or declared type of an environment variable value.
type Type string

const (
	TypeString  Type = "string"
	TypeBool    Type = "bool"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeURL     Type = "url"
	TypeUnknown Type = "unknown"
)

// Supported returns all known type names.
func Supported() []string {
	return []string{string(TypeString), string(TypeBool), string(TypeInt), string(TypeFloat), string(TypeURL)}
}

// IsSupported reports whether t is a recognised type.
func IsSupported(t string) bool {
	for _, s := range Supported() {
		if s == t {
			return true
		}
	}
	return false
}

// Infer guesses the Type of a raw string value.
func Infer(value string) Type {
	if value == "" {
		return TypeString
	}
	lower := strings.ToLower(value)
	if lower == "true" || lower == "false" {
		return TypeBool
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return TypeFloat
	}
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return TypeURL
	}
	return TypeString
}

// InferAll returns a map of key → Type for every entry in vars.
func InferAll(vars map[string]string) map[string]Type {
	out := make(map[string]Type, len(vars))
	for k, v := range vars {
		out[k] = Infer(v)
	}
	return out
}

// Validate checks whether value is compatible with the declared type t.
// It returns a descriptive error when the value cannot be parsed as t.
func Validate(key, value string, t Type) error {
	switch t {
	case TypeBool:
		l := strings.ToLower(value)
		if l != "true" && l != "false" {
			return fmt.Errorf("%s: %q is not a valid bool (expected true/false)", key, value)
		}
	case TypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("%s: %q is not a valid int", key, value)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("%s: %q is not a valid float", key, value)
		}
	case TypeURL:
		l := strings.ToLower(value)
		if !strings.HasPrefix(l, "http://") && !strings.HasPrefix(l, "https://") {
			return fmt.Errorf("%s: %q is not a valid URL (must start with http:// or https://)", key, value)
		}
	}
	return nil
}
