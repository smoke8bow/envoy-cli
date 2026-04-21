package masked_export_test

import (
	"strings"
	"testing"

	"github.com/your-org/envoy-cli/internal/masked_export"
)

func TestNewExporterValid(t *testing.T) {
	for _, f := range []masked_export.Format{
		masked_export.FormatDotenv,
		masked_export.FormatShell,
		masked_export.FormatJSON,
	} {
		_, err := masked_export.NewExporter(f, nil, nil)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNewExporterInvalid(t *testing.T) {
	_, err := masked_export.NewExporter("xml", nil, nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestRenderDotenvMasksSensitive(t *testing.T) {
	e, _ := masked_export.NewExporter(masked_export.FormatDotenv, nil, nil)
	vars := map[string]string{
		"API_KEY": "secret123",
		"APP_NAME": "myapp",
	}
	out, err := e.Render(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "secret123") {
		t.Errorf("expected sensitive value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME to be present in output")
	}
}

func TestRenderShellMasksSensitive(t *testing.T) {
	e, _ := masked_export.NewExporter(masked_export.FormatShell, nil, nil)
	vars := map[string]string{"SECRET_TOKEN": "topsecret"}
	out, err := e.Render(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "topsecret") {
		t.Errorf("expected value to be masked in shell output, got: %s", out)
	}
	if !strings.Contains(out, "export SECRET_TOKEN") {
		t.Errorf("expected export statement in shell output")
	}
}

func TestRenderRevealedKeyNotMasked(t *testing.T) {
	e, _ := masked_export.NewExporter(masked_export.FormatDotenv, nil, []string{"API_KEY"})
	vars := map[string]string{"API_KEY": "visible_secret"}
	out, err := e.Render(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "visible_secret") {
		t.Errorf("expected revealed key value to appear unmasked, got: %s", out)
	}
}

func TestRenderJSONContainsKeys(t *testing.T) {
	e, _ := masked_export.NewExporter(masked_export.FormatJSON, nil, nil)
	vars := map[string]string{"DB_PASS": "hunter2", "REGION": "us-east-1"}
	out, err := e.Render(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS key in JSON output, got: %s", out)
	}
	if strings.Contains(out, "hunter2") {
		t.Errorf("expected sensitive value to be masked in JSON output, got: %s", out)
	}
}

func TestRenderCustomPatternMasks(t *testing.T) {
	e, _ := masked_export.NewExporter(masked_export.FormatDotenv, []string{"CUSTOM_.*"}, nil)
	vars := map[string]string{"CUSTOM_FIELD": "shouldbehidden", "OTHER": "visible"}
	out, err := e.Render(vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "shouldbehidden") {
		t.Errorf("expected custom pattern match to be masked, got: %s", out)
	}
	if !strings.Contains(out, "visible") {
		t.Errorf("expected non-matching key value to be visible, got: %s", out)
	}
}
