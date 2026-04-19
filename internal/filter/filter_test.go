package filter

import (
	"testing"
)

var sampleVars = map[string]string{
	"AWS_ACCESS_KEY":  "abc",
	"AWS_SECRET_KEY":  "xyz",
	"DB_HOST":         "localhost",
	"DB_PORT":         "5432",
	"APP_DEBUG":       "true",
	"APP_SECRET_KEY":  "s3cr3t",
}

func TestFilterByPrefix(t *testing.T) {
	res := Filter(sampleVars, Option{Prefix: "AWS_"})
	if len(res.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(res.Matched))
	}
	if _, ok := res.Matched["AWS_ACCESS_KEY"]; !ok {
		t.Error("expected AWS_ACCESS_KEY in matched")
	}
}

func TestFilterBySuffix(t *testing.T) {
	res := Filter(sampleVars, Option{Suffix: "_KEY"})
	if len(res.Matched) != 3 {
		t.Fatalf("expected 3 matched, got %d", len(res.Matched))
	}
}

func TestFilterByContains(t *testing.T) {
	res := Filter(sampleVars, Option{Contains: "SECRET"})
	if len(res.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(res.Matched))
	}
}

func TestFilterByExactKeys(t *testing.T) {
	res := Filter(sampleVars, Option{ExactKeys: []string{"DB_HOST", "DB_PORT"}})
	if len(res.Matched) != 2 {
		t.Fatalf("expected 2 matched, got %d", len(res.Matched))
	}
	if len(res.Excluded) != 4 {
		t.Fatalf("expected 4 excluded, got %d", len(res.Excluded))
	}
}

func TestFilterCombinedPrefixAndSuffix(t *testing.T) {
	res := Filter(sampleVars, Option{Prefix: "APP_", Suffix: "_KEY"})
	if len(res.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(res.Matched))
	}
	if _, ok := res.Matched["APP_SECRET_KEY"]; !ok {
		t.Error("expected APP_SECRET_KEY in matched")
	}
}

func TestFilterNoOptions(t *testing.T) {
	res := Filter(sampleVars, Option{})
	if len(res.Matched) != len(sampleVars) {
		t.Fatalf("expected all keys matched, got %d", len(res.Matched))
	}
}

func TestFilterExactKeysIgnoresOtherOptions(t *testing.T) {
	res := Filter(sampleVars, Option{
		Prefix:    "AWS_",
		ExactKeys: []string{"DB_HOST"},
	})
	if len(res.Matched) != 1 {
		t.Fatalf("expected 1 matched, got %d", len(res.Matched))
	}
	if _, ok := res.Matched["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in matched")
	}
}
