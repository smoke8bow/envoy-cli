package lint

import (
	"testing"
)

func TestNoIssues(t *testing.T) {
	l := NewLinter()
	vars := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	issues := l.Run(vars)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRuleNoEmptyValues(t *testing.T) {
	vars := map[string]string{
		"KEY": "",
		"OTHER": "value",
	}
	issues := RuleNoEmptyValues(vars)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "KEY" {
		t.Errorf("expected KEY, got %s", issues[0].Key)
	}
}

func TestRuleNoWhitespaceInKeys(t *testing.T) {
	vars := map[string]string{
		"BAD KEY": "value",
		"GOOD_KEY": "value",
	}
	issues := RuleNoWhitespaceInKeys(vars)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "BAD KEY" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestRuleNoLowercaseKeys(t *testing.T) {
	vars := map[string]string{
		"lowercase": "value",
		"MixedCase": "value",
		"UPPERCASE": "value",
	}
	issues := RuleNoLowercaseKeys(vars)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
}

func TestIssueString(t *testing.T) {
	i := Issue{Key: "FOO", Message: "some problem"}
	got := i.String()
	want := "FOO: some problem"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMultipleRulesAggregated(t *testing.T) {
	l := NewLinter()
	vars := map[string]string{
		"bad key": "",
	}
	issues := l.Run(vars)
	// empty value + whitespace in key + not uppercase = 3 issues
	if len(issues) != 3 {
		t.Fatalf("expected 3 issues, got %d: %v", len(issues), issues)
	}
}
