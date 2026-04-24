package extract

import (
	"testing"
)

var sampleVars = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"APP_ENV":  "production",
	"SECRET":  "s3cr3t",
}

func TestExtractAllWhenNoKeys(t *testing.T) {
	res, err := Extract(sampleVars, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != len(sampleVars) {
		t.Fatalf("expected %d vars, got %d", len(sampleVars), len(res.Vars))
	}
}

func TestExtractDoesNotMutateSource(t *testing.T) {
	res, _ := Extract(sampleVars, Options{})
	res.Vars["INJECTED"] = "value"
	if _, ok := sampleVars["INJECTED"]; ok {
		t.Fatal("source map was mutated")
	}
}

func TestExtractSelectedKeys(t *testing.T) {
	res, err := Extract(sampleVars, Options{Keys: []string{"DB_HOST", "APP_ENV"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(res.Vars))
	}
	if res.Vars["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", res.Vars["DB_HOST"])
	}
}

func TestExtractMissingKeysSoft(t *testing.T) {
	res, err := Extract(sampleVars, Options{Keys: []string{"DB_HOST", "MISSING_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in missing list, got %v", res.Missing)
	}
	if _, ok := res.Vars["MISSING_KEY"]; ok {
		t.Error("missing key should not appear in Vars")
	}
}

func TestExtractMissingKeysFailOnMissing(t *testing.T) {
	_, err := Extract(sampleVars, Options{
		Keys:          []string{"DB_HOST", "NO_SUCH_KEY"},
		FailOnMissing: true,
	})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestExtractEmptySource(t *testing.T) {
	res, err := Extract(map[string]string{}, Options{Keys: []string{"KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 missing key, got %d", len(res.Missing))
	}
}

func TestExtractDuplicateKeys(t *testing.T) {
	res, err := Extract(sampleVars, Options{Keys: []string{"DB_HOST", "DB_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 1 {
		t.Errorf("expected 1 var for duplicate keys, got %d", len(res.Vars))
	}
}
