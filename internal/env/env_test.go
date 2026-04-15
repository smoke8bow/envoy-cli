package env

import (
	"os"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	vars := map[string]string{
		"ENVOY_TEST_FOO": "bar",
		"ENVOY_TEST_BAZ": "qux",
	}

	if err := Apply(vars); err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	for key, want := range vars {
		got := os.Getenv(key)
		if got != want {
			t.Errorf("os.Getenv(%q) = %q, want %q", key, got, want)
		}
	}

	// Cleanup
	_ = Unset([]string{"ENVOY_TEST_FOO", "ENVOY_TEST_BAZ"})
}

func TestExport(t *testing.T) {
	vars := map[string]string{
		"MY_VAR": "hello world",
		"QUOTED": `say "hi"`,
	}

	output := Export(vars)

	if !strings.Contains(output, `export MY_VAR="hello world"`) {
		t.Errorf("Export() missing MY_VAR line, got:\n%s", output)
	}
	if !strings.Contains(output, `export QUOTED="say \"hi\""`) {
		t.Errorf("Export() missing properly escaped QUOTED line, got:\n%s", output)
	}
}

func TestUnset(t *testing.T) {
	_ = os.Setenv("ENVOY_UNSET_TEST", "value")

	if err := Unset([]string{"ENVOY_UNSET_TEST"}); err != nil {
		t.Fatalf("Unset() unexpected error: %v", err)
	}

	if val := os.Getenv("ENVOY_UNSET_TEST"); val != "" {
		t.Errorf("expected ENVOY_UNSET_TEST to be unset, got %q", val)
	}
}

func TestSnapshot(t *testing.T) {
	_ = os.Setenv("ENVOY_SNAP_A", "alpha")
	_ = os.Setenv("ENVOY_SNAP_B", "beta")
	defer func() {
		_ = Unset([]string{"ENVOY_SNAP_A", "ENVOY_SNAP_B"})
	}()

	snap := Snapshot([]string{"ENVOY_SNAP_A", "ENVOY_SNAP_B", "ENVOY_SNAP_MISSING"})

	if snap["ENVOY_SNAP_A"] != "alpha" {
		t.Errorf("expected ENVOY_SNAP_A=alpha, got %q", snap["ENVOY_SNAP_A"])
	}
	if snap["ENVOY_SNAP_B"] != "beta" {
		t.Errorf("expected ENVOY_SNAP_B=beta, got %q", snap["ENVOY_SNAP_B"])
	}
	if _, ok := snap["ENVOY_SNAP_MISSING"]; ok {
		t.Error("expected ENVOY_SNAP_MISSING to be absent from snapshot")
	}
}
