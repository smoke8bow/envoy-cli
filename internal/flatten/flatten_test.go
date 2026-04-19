package flatten

import (
	"testing"
)

func TestFlattenSimple(t *testing.T) {
	input := map[string]any{
		"host": "localhost",
		"port": "5432",
	}
	out := Flatten(input, DefaultOptions())
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestFlattenNested(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	out := Flatten(input, DefaultOptions())
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlattenDeepNested(t *testing.T) {
	input := map[string]any{
		"app": map[string]any{
			"db": map[string]any{
				"pass": "secret",
			},
		},
	}
	out := Flatten(input, DefaultOptions())
	if out["APP_DB_PASS"] != "secret" {
		t.Errorf("expected APP_DB_PASS=secret, got %q", out["APP_DB_PASS"])
	}
}

func TestFlattenWithPrefix(t *testing.T) {
	input := map[string]any{"key": "val"}
	opts := Options{Separator: "_", Uppercase: true, Prefix: "myapp"}
	out := Flatten(input, opts)
	if out["MYAPP_KEY"] != "val" {
		t.Errorf("expected MYAPP_KEY=val, got %q", out["MYAPP_KEY"])
	}
}

func TestFlattenNonStringValue(t *testing.T) {
	input := map[string]any{"count": 42}
	out := Flatten(input, DefaultOptions())
	if out["COUNT"] != "42" {
		t.Errorf("expected COUNT=42, got %q", out["COUNT"])
	}
}

func TestFlattenLowercaseKeys(t *testing.T) {
	input := map[string]any{"Key": "value"}
	opts := Options{Separator: "_", Uppercase: false}
	out := Flatten(input, opts)
	if out["Key"] != "value" {
		t.Errorf("expected Key=value, got %q", out["Key"])
	}
}

func TestFlattenCustomSeparator(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"host": "localhost"},
	}
	opts := Options{Separator: ".", Uppercase: false}
	out := Flatten(input, opts)
	if out["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %q", out["db.host"])
	}
}
