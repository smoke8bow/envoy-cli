package envfence_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envfence"
)

func TestNewInvalidMode(t *testing.T) {
	_, err := envfence.New("unknown", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewEmptyKeys(t *testing.T) {
	_, err := envfence.New(envfence.ModeAllow, nil)
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestCheckAllowlistNoViolations(t *testing.T) {
	f, err := envfence.New(envfence.ModeAllow, []string{"HOST", "PORT"})
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{"HOST": "localhost", "PORT": "8080"}
	violations := f.Check(vars)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckAllowlistViolation(t *testing.T) {
	f, err := envfence.New(envfence.ModeAllow, []string{"HOST"})
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{"HOST": "localhost", "SECRET": "abc"}
	violations := f.Check(vars)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET" {
		t.Errorf("expected violation for SECRET, got %s", violations[0].Key)
	}
}

func TestCheckDenylistNoViolations(t *testing.T) {
	f, err := envfence.New(envfence.ModeDeny, []string{"SECRET"})
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{"HOST": "localhost", "PORT": "8080"}
	violations := f.Check(vars)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckDenylistViolation(t *testing.T) {
	f, err := envfence.New(envfence.ModeDeny, []string{"SECRET"})
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{"HOST": "localhost", "SECRET": "abc"}
	violations := f.Check(vars)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestFilterAllowlist(t *testing.T) {
	f, _ := envfence.New(envfence.ModeAllow, []string{"HOST", "PORT"})
	vars := map[string]string{"HOST": "localhost", "PORT": "8080", "SECRET": "xyz"}
	out := f.Filter(vars)
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should have been filtered out")
	}
	if out["HOST"] != "localhost" {
		t.Error("HOST should be present")
	}
}

func TestFilterDenylist(t *testing.T) {
	f, _ := envfence.New(envfence.ModeDeny, []string{"SECRET"})
	vars := map[string]string{"HOST": "localhost", "SECRET": "xyz"}
	out := f.Filter(vars)
	if _, ok := out["SECRET"]; ok {
		t.Error("SECRET should have been filtered out")
	}
	if out["HOST"] != "localhost" {
		t.Error("HOST should be present")
	}
}

func TestFilterDoesNotMutateInput(t *testing.T) {
	f, _ := envfence.New(envfence.ModeAllow, []string{"HOST"})
	vars := map[string]string{"HOST": "localhost", "SECRET": "xyz"}
	f.Filter(vars)
	if _, ok := vars["SECRET"]; !ok {
		t.Error("Filter must not mutate the input map")
	}
}
