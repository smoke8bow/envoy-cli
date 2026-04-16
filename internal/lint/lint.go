package lint

import (
	"fmt"
	"strings"
)

// Issue represents a single lint finding.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Rule is a function that inspects env vars and returns issues.
type Rule func(vars map[string]string) []Issue

// Linter runs a set of rules against a profile's env vars.
type Linter struct {
	rules []Rule
}

// NewLinter returns a Linter with the default rule set.
func NewLinter() *Linter {
	return &Linter{
		rules: []Rule{
			RuleNoEmptyValues,
			RuleNoWhitespaceInKeys,
			RuleNoLowercaseKeys,
		},
	}
}

// Run executes all rules and returns aggregated issues.
func (l *Linter) Run(vars map[string]string) []Issue {
	var issues []Issue
	for _, rule := range l.rules {
		issues = append(issues, rule(vars)...)
	}
	return issues
}

// RuleNoEmptyValues flags keys with empty string values.
func RuleNoEmptyValues(vars map[string]string) []Issue {
	var issues []Issue
	for k, v := range vars {
		if strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Key: k, Message: "value is empty"})
		}
	}
	return issues
}

// RuleNoWhitespaceInKeys flags keys containing whitespace.
func RuleNoWhitespaceInKeys(vars map[string]string) []Issue {
	var issues []Issue
	for k := range vars {
		if strings.ContainsAny(k, " \t") {
			issues = append(issues, Issue{Key: k, Message: "key contains whitespace"})
		}
	}
	return issues
}

// RuleNoLowercaseKeys flags keys that are not fully uppercase.
func RuleNoLowercaseKeys(vars map[string]string) []Issue {
	var issues []Issue
	for k := range vars {
		if k != strings.ToUpper(k) {
			issues = append(issues, Issue{Key: k, Message: "key is not uppercase"})
		}
	}
	return issues
}
