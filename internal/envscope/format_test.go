package envscope

import (
	"strings"
	"testing"
)

func TestFormatWithHeader(t *testing.T) {
	sv := &ScopedView{
		Profile: "prod",
		Scope:   "APP_",
		Vars:    map[string]string{"HOST": "localhost", "PORT": "8080"},
	}
	out := Format(sv, DefaultFormatOptions())
	if !strings.Contains(out, "# profile=prod scope=APP_") {
		t.Errorf("expected header line, got:\n%s", out)
	}
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", out)
	}
}

func TestFormatWithoutHeader(t *testing.T) {
	sv := &ScopedView{
		Profile: "dev",
		Scope:   "DB_",
		Vars:    map[string]string{"HOST": "db"},
	}
	out := Format(sv, FormatOptions{ShowProfile: false})
	if strings.Contains(out, "#") {
		t.Errorf("expected no header, got:\n%s", out)
	}
	if !strings.Contains(out, "HOST=db") {
		t.Errorf("expected HOST=db, got:\n%s", out)
	}
}

func TestFormatSortedOutput(t *testing.T) {
	sv := &ScopedView{
		Profile: "staging",
		Scope:   "SVC_",
		Vars:    map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"},
	}
	out := Format(sv, FormatOptions{ShowProfile: false})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to start with A_KEY, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to start with Z_KEY, got %q", lines[2])
	}
}

func TestFormatEmptyVars(t *testing.T) {
	sv := &ScopedView{
		Profile: "empty",
		Scope:   "NONE_",
		Vars:    map[string]string{},
	}
	out := Format(sv, DefaultFormatOptions())
	if !strings.Contains(out, "# profile=empty") {
		t.Errorf("expected header even for empty vars, got:\n%s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
