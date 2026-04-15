package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrEmptyName     = errors.New("name cannot be empty")
	ErrInvalidName   = errors.New("name contains invalid characters")
	ErrEmptyKey      = errors.New("environment variable key cannot be empty")
	ErrInvalidKey    = errors.New("environment variable key contains invalid characters")
	ErrReservedKey   = errors.New("environment variable key is reserved")
)

// namePattern allows alphanumerics, hyphens, and underscores.
var namePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// keyPattern follows POSIX env var naming conventions.
var keyPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// reservedKeys are keys that should not be overwritten by user profiles.
var reservedKeys = map[string]struct{}{
	"PATH":  {},
	"HOME":  {},
	"USER":  {},
	"SHELL": {},
	"PWD":   {},
}

// ProfileName validates a profile name.
func ProfileName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	if !namePattern.MatchString(name) {
		return fmt.Errorf("%w: %q", ErrInvalidName, name)
	}
	return nil
}

// EnvKey validates a single environment variable key.
func EnvKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("%w: %q", ErrInvalidKey, key)
	}
	if _, reserved := reservedKeys[strings.ToUpper(key)]; reserved {
		return fmt.Errorf("%w: %q", ErrReservedKey, key)
	}
	return nil
}

// EnvVars validates a map of environment variable key-value pairs.
func EnvVars(vars map[string]string) error {
	for k := range vars {
		if err := EnvKey(k); err != nil {
			return err
		}
	}
	return nil
}
