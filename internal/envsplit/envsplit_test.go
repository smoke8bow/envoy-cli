package envsplit_test

import (
	"testing"

	"envoy-cli/internal/envsplit"
)

func TestSplitNoRules(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "DB_URL": "postgres://"}
	res, err := envsplit.Split(src, nil, envsplit.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets) != 0 {
		t.Errorf("expected no buckets, got %d", len(res.Buckets))
	}
	if len(res.Remainder) != 2 {
		t.Errorf("expected 2 remainder keys, got %d", len(res.Remainder))
	}
}

func TestSplitBasic(t *testing.T) {
	src := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://",
		"OTHER":    "val",
	}
	rules := []envsplit.Rule{
		{Prefix: "APP_", Bucket: "app"},
		{Prefix: "DB_", Bucket: "db"},
	}
	res, err := envsplit.Split(src, rules, envsplit.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["app"]) != 2 {
		t.Errorf("app bucket: want 2, got %d", len(res.Buckets["app"]))
	}
	if res.Buckets["app"]["APP_HOST"] != "localhost" {
		t.Errorf("APP_HOST not preserved")
	}
	if len(res.Buckets["db"]) != 1 {
		t.Errorf("db bucket: want 1, got %d", len(res.Buckets["db"]))
	}
	if len(res.Remainder) != 1 {
		t.Errorf("remainder: want 1, got %d", len(res.Remainder))
	}
}

func TestSplitStripPrefix(t *testing.T) {
	src := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	rules := []envsplit.Rule{{Prefix: "APP_", Bucket: "app"}}
	opts := envsplit.Options{StripPrefix: true}
	res, err := envsplit.Split(src, rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Buckets["app"]["HOST"]; !ok {
		t.Error("expected stripped key HOST in app bucket")
	}
	if _, ok := res.Buckets["app"]["PORT"]; !ok {
		t.Error("expected stripped key PORT in app bucket")
	}
}

func TestSplitFirstRuleWins(t *testing.T) {
	src := map[string]string{"APP_DB_URL": "postgres://"}
	rules := []envsplit.Rule{
		{Prefix: "APP_", Bucket: "app"},
		{Prefix: "APP_DB_", Bucket: "db"},
	}
	res, err := envsplit.Split(src, rules, envsplit.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets["app"]) != 1 {
		t.Error("expected APP_DB_URL in app bucket (first rule wins)")
	}
	if len(res.Buckets["db"]) != 0 {
		t.Error("expected db bucket to be empty")
	}
}

func TestSplitEmptyPrefixError(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	rules := []envsplit.Rule{{Prefix: "", Bucket: "x"}}
	_, err := envsplit.Split(src, rules, envsplit.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestSplitEmptyBucketError(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	rules := []envsplit.Rule{{Prefix: "KEY", Bucket: ""}}
	_, err := envsplit.Split(src, rules, envsplit.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty bucket name")
	}
}
