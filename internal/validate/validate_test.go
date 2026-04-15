package validate_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envoy-cli/internal/validate"
)

func TestProfileNameValid(t *testing.T) {
	valid := []string{"dev", "prod", "my-env", "env_1", "ABC", "a1b2-c3"}
	for _, name := range valid {
		if err := validate.ProfileName(name); err != nil {
			t.Errorf("expected %q to be valid, got: %v", name, err)
		}
	}
}

func TestProfileNameInvalid(t *testing.T) {
	invalid := []string{"", "  ", "my env", "env@1", "env.name", "env/path"}
	for _, name := range invalid {
		if err := validate.ProfileName(name); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", name)
		}
	}
}

func TestProfileNameEmptyError(t *testing.T) {
	err := validate.ProfileName("")
	if !errors.Is(err, validate.ErrEmptyName) {
		t.Errorf("expected ErrEmptyName, got: %v", err)
	}
}

func TestEnvKeyValid(t *testing.T) {
	valid := []string{"FOO", "BAR_BAZ", "_PRIVATE", "MY_VAR_1"}
	for _, key := range valid {
		if err := validate.EnvKey(key); err != nil {
			t.Errorf("expected %q to be valid, got: %v", key, err)
		}
	}
}

func TestEnvKeyInvalid(t *testing.T) {
	invalid := []string{"", "1STARTS_WITH_NUM", "HAS SPACE", "HAS-HYPHEN", "DOT.KEY"}
	for _, key := range invalid {
		if err := validate.EnvKey(key); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", key)
		}
	}
}

func TestEnvKeyReserved(t *testing.T) {
	reserved := []string{"PATH", "HOME", "USER", "SHELL", "PWD"}
	for _, key := range reserved {
		err := validate.EnvKey(key)
		if !errors.Is(err, validate.ErrReservedKey) {
			t.Errorf("expected ErrReservedKey for %q, got: %v", key, err)
		}
	}
}

func TestEnvVarsValid(t *testing.T) {
	vars := map[string]string{
		"FOO":     "bar",
		"MY_VAR":  "value",
		"_SECRET": "token",
	}
	if err := validate.EnvVars(vars); err != nil {
		t.Errorf("expected valid vars, got: %v", err)
	}
}

func TestEnvVarsInvalidKey(t *testing.T) {
	vars := map[string]string{
		"GOOD_KEY": "ok",
		"bad-key":  "fail",
	}
	if err := validate.EnvVars(vars); err == nil {
		t.Error("expected error for invalid key in map, got nil")
	}
}
