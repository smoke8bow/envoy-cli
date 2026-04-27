package envtrim_test

import (
	"testing"

	"github.com/nicholasgasior/envoy-cli/internal/envtrim"
)

func TestDefaultOptionsTrimsAll(t *testing.T) {
	opts := envtrim.DefaultOptions()
	if !opts.TrimKeys || !opts.TrimValues {
		t.Fatal("expected both TrimKeys and TrimValues to be true")
	}
}

func TestTrimValuesOnly(t *testing.T) {
	src := map[string]string{
		"KEY": "  hello  ",
		"OTHER": "world",
	}
	res, err := envtrim.Trim(src, envtrim.Options{TrimKeys: false, TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Vars["KEY"]; got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
	if res.Changes != 1 {
		t.Errorf("expected 1 change, got %d", res.Changes)
	}
}

func TestTrimKeysOnly(t *testing.T) {
	src := map[string]string{
		" SPACED ": "value",
	}
	res, err := envtrim.Trim(src, envtrim.Options{TrimKeys: true, TrimValues: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Vars["SPACED"]; !ok {
		t.Error("expected trimmed key 'SPACED' to exist")
	}
	if res.Changes != 1 {
		t.Errorf("expected 1 change, got %d", res.Changes)
	}
}

func TestTrimBothKeysAndValues(t *testing.T) {
	src := map[string]string{
		" FOO ": "  bar  ",
		"CLEAN": "value",
	}
	res, err := envtrim.Trim(src, envtrim.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Vars["FOO"]; got != "bar" {
		t.Errorf("expected 'bar', got %q", got)
	}
	if res.Changes != 1 {
		t.Errorf("expected 1 change, got %d", res.Changes)
	}
}

func TestTrimDoesNotMutateSource(t *testing.T) {
	src := map[string]string{"KEY": "  val  "}
	_, err := envtrim.Trim(src, envtrim.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["KEY"] != "  val  " {
		t.Error("source map was mutated")
	}
}

func TestTrimNeitherReturnsError(t *testing.T) {
	_, err := envtrim.Trim(map[string]string{}, envtrim.Options{})
	if err == nil {
		t.Fatal("expected error when neither TrimKeys nor TrimValues is set")
	}
}

func TestTrimNoChanges(t *testing.T) {
	src := map[string]string{"KEY": "value", "OTHER": "data"}
	res, err := envtrim.Trim(src, envtrim.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changes != 0 {
		t.Errorf("expected 0 changes, got %d", res.Changes)
	}
}
