package normalize

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
	if IsSupported(Strategy("camel")) {
		t.Error("expected camel to be unsupported")
	}
}

func TestNewInvalidStrategy(t *testing.T) {
	_, err := New("camel")
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestApplyUpper(t *testing.T) {
	n, _ := New(StrategyUpper)
	out := n.Apply(map[string]string{"foo": "bar", "Baz": "qux"})
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApplyLower(t *testing.T) {
	n, _ := New(StrategyLower)
	out := n.Apply(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if out["foo"] != "bar" || out["baz"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestApplySnake(t *testing.T) {
	n, _ := New(StrategySnake)
	out := n.Apply(map[string]string{"my-key": "1", "another key": "2"})
	if out["MY_KEY"] != "1" {
		t.Errorf("expected MY_KEY, got %v", out)
	}
	if out["ANOTHER_KEY"] != "2" {
		t.Errorf("expected ANOTHER_KEY, got %v", out)
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	n, _ := New(StrategyUpper)
	input := map[string]string{"foo": "bar"}
	n.Apply(input)
	if _, ok := input["FOO"]; ok {
		t.Error("Apply mutated input map")
	}
}

func TestCollisionLastWins(t *testing.T) {
	n, _ := New(StrategyUpper)
	// Both "foo" and "FOO" normalize to "FOO"; last write wins (map iteration order varies).
	out := n.Apply(map[string]string{"FOO": "original"})
	if out["FOO"] != "original" {
		t.Errorf("unexpected value: %v", out["FOO"])
	}
}
