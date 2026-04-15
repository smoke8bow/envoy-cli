package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewExporterValid(t *testing.T) {
	for _, f := range []Format{FormatJSON, FormatDotenv, FormatShell} {
		_, err := NewExporter(f)
		if err != nil {
			t.Errorf("expected no error for format %s, got %v", f, err)
		}
	}
}

func TestNewExporterInvalid(t *testing.T) {
	_, err := NewExporter("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestRenderJSON(t *testing.T) {
	ex, _ := NewExporter(FormatJSON)
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := ex.Render("myprofile", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var ep ExportedProfile
	if err := json.Unmarshal([]byte(out), &ep); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if ep.Name != "myprofile" {
		t.Errorf("expected name myprofile, got %s", ep.Name)
	}
	if ep.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
}

func TestRenderDotenv(t *testing.T) {
	ex, _ := NewExporter(FormatDotenv)
	vars := map[string]string{"KEY": "value"}
	out, err := ex.Render("p", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY=value") {
		t.Errorf("expected KEY=value in dotenv output, got: %s", out)
	}
}

func TestRenderShell(t *testing.T) {
	ex, _ := NewExporter(FormatShell)
	vars := map[string]string{"MY_VAR": "hello world"}
	out, err := ex.Render("p", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export MY_VAR=") {
		t.Errorf("expected export statement in shell output, got: %s", out)
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "subdir", "profile.env")
	ex, _ := NewExporter(FormatDotenv)
	vars := map[string]string{"A": "1"}
	if err := ex.WriteFile(dest, "test", vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if !strings.Contains(string(data), "A=1") {
		t.Errorf("expected A=1 in file, got: %s", string(data))
	}
}
