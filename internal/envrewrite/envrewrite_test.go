package envrewrite_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envrewrite"
)

// newRewriter returns a Rewriter with the given rules applied.
func newRewriter(rules []envrewrite.Rule) *envrewrite.Rewriter {
	return envrewrite.New(rules)
}

func TestApplyNoRules(t *testing.T) {
	r := newRewriter(nil)
	input := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(input) {
		t.Fatalf("expected %d keys, got %d", len(input), len(out))
	}
	for k, v := range input {
		if out[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, out[k])
		}
	}
}

func TestRenameKey(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "FOO", NewKey: "BAR"},
	}
	r := newRewriter(rules)
	input := map[string]string{"FOO": "hello", "OTHER": "world"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["FOO"]; ok {
		t.Error("old key FOO should not exist after rename")
	}
	if out["BAR"] != "hello" {
		t.Errorf("expected BAR=hello, got %q", out["BAR"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected OTHER=world, got %q", out["OTHER"])
	}
}

func TestRenameKeyNotFound(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "MISSING", NewKey: "NEW_KEY"},
	}
	r := newRewriter(rules)
	input := map[string]string{"FOO": "bar"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// MISSING key does not exist; output should be unchanged
	if len(out) != 1 || out["FOO"] != "bar" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestRewriteValue(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "DB_URL", NewKey: "DB_URL", ValueTemplate: "postgres://localhost/{{.Value}}"},
	}
	r := newRewriter(rules)
	input := map[string]string{"DB_URL": "mydb"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://localhost/mydb"
	if out["DB_URL"] != want {
		t.Errorf("expected %q, got %q", want, out["DB_URL"])
	}
}

func TestRenameAndRewriteValue(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "OLD_HOST", NewKey: "APP_HOST", ValueTemplate: "https://{{.Value}}"},
	}
	r := newRewriter(rules)
	input := map[string]string{"OLD_HOST": "example.com", "PORT": "8080"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_HOST"]; ok {
		t.Error("OLD_HOST should have been renamed")
	}
	if out["APP_HOST"] != "https://example.com" {
		t.Errorf("expected APP_HOST=https://example.com, got %q", out["APP_HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "FOO", NewKey: "BAR"},
	}
	r := newRewriter(rules)
	input := map[string]string{"FOO": "original"}
	_, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input["FOO"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestMultipleRules(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "A", NewKey: "X"},
		{OldKey: "B", NewKey: "Y"},
	}
	r := newRewriter(rules)
	input := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "1" {
		t.Errorf("expected X=1, got %q", out["X"])
	}
	if out["Y"] != "2" {
		t.Errorf("expected Y=2, got %q", out["Y"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %q", out["C"])
	}
	if _, ok := out["A"]; ok {
		t.Error("A should have been renamed to X")
	}
	if _, ok := out["B"]; ok {
		t.Error("B should have been renamed to Y")
	}
}

func TestInvalidValueTemplate(t *testing.T) {
	rules := []envrewrite.Rule{
		{OldKey: "FOO", NewKey: "FOO", ValueTemplate: "{{invalid"},
	}
	// New should surface the error at construction time or Apply time
	r := newRewriter(rules)
	_, err := r.Apply(map[string]string{"FOO": "bar"})
	if err == nil {
		t.Error("expected error for invalid template, got nil")
	}
}
