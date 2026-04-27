package envgroup

import (
	"strings"
	"testing"
)

func TestGroupByBasic(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
		"APP_PORT": "8080",
		"SINGLE":  "value",
	}
	r := GroupBy(vars, DefaultOptions())
	if len(r.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(r.Groups))
	}
	if r.Groups[0].Prefix != "APP" {
		t.Errorf("expected APP first, got %s", r.Groups[0].Prefix)
	}
	if _, ok := r.Ungrouped["SINGLE"]; !ok {
		t.Error("SINGLE should be ungrouped")
	}
}

func TestGroupByMinSize(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"DB_NAME": "mydb",
	}
	opts := DefaultOptions()
	opts.MinSize = 3
	r := GroupBy(vars, opts)
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(r.Groups))
	}
	if r.Groups[0].Prefix != "DB" {
		t.Errorf("expected DB group")
	}
}

func TestGroupByMinSizeNotMet(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"OTHER":   "val",
	}
	opts := DefaultOptions()
	opts.MinSize = 3
	r := GroupBy(vars, opts)
	if len(r.Groups) != 0 {
		t.Fatalf("expected no groups, got %d", len(r.Groups))
	}
	if len(r.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped, got %d", len(r.Ungrouped))
	}
}

func TestGroupByStripPrefix(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	opts := DefaultOptions()
	opts.StripPrefix = true
	r := GroupBy(vars, opts)
	if len(r.Groups) != 1 {
		t.Fatalf("expected 1 group")
	}
	g := r.Groups[0]
	if _, ok := g.Vars["HOST"]; !ok {
		t.Error("expected key HOST after strip")
	}
	if _, ok := g.Vars["PORT"]; !ok {
		t.Error("expected key PORT after strip")
	}
}

func TestGroupByCustomSeparator(t *testing.T) {
	vars := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
		"app.env": "dev",
		"app.port": "9000",
	}
	opts := Options{MinSize: 2, Separator: "."}
	r := GroupBy(vars, opts)
	if len(r.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(r.Groups))
	}
}

func TestGroupByEmpty(t *testing.T) {
	r := GroupBy(map[string]string{}, DefaultOptions())
	if len(r.Groups) != 0 {
		t.Error("expected no groups for empty input")
	}
}

func TestFormatContainsPrefix(t *testing.T) {
	vars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	r := GroupBy(vars, DefaultOptions())
	out := Format(r)
	if !strings.Contains(out, "[DB]") {
		t.Errorf("expected [DB] in output, got:\n%s", out)
	}
}

func TestFormatUngrouped(t *testing.T) {
	vars := map[string]string{
		"LONE": "val",
	}
	r := GroupBy(vars, DefaultOptions())
	out := Format(r)
	if !strings.Contains(out, "[ungrouped]") {
		t.Errorf("expected [ungrouped] section, got:\n%s", out)
	}
}
