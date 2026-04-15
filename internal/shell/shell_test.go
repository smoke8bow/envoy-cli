package shell

import (
	"strings"
	"testing"
)

func TestNewExporterDefaults(t *testing.T) {
	e := NewExporter("unknown")
	if e.Shell != Bash {
		t.Errorf("expected bash, got %s", e.Shell)
	}
}

func TestNewExporterZsh(t *testing.T) {
	e := NewExporter("zsh")
	if e.Shell != Zsh {
		t.Errorf("expected zsh, got %s", e.Shell)
	}
}

func TestNewExporterFish(t *testing.T) {
	e := NewExporter("fish")
	if e.Shell != Fish {
		t.Errorf("expected fish, got %s", e.Shell)
	}
}

func TestExportStatementsBash(t *testing.T) {
	e := NewExporter("bash")
	env := map[string]string{"FOO": "bar"}
	lines := e.ExportStatements(env)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "export FOO=") {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestExportStatementsFish(t *testing.T) {
	e := NewExporter("fish")
	env := map[string]string{"FOO": "bar"}
	lines := e.ExportStatements(env)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "set -x FOO") {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestUnsetStatementsBash(t *testing.T) {
	e := NewExporter("bash")
	lines := e.UnsetStatements([]string{"FOO", "BAR"})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, l := range lines {
		if !strings.HasPrefix(l, "unset ") {
			t.Errorf("unexpected line: %s", l)
		}
	}
}

func TestUnsetStatementsFish(t *testing.T) {
	e := NewExporter("fish")
	lines := e.UnsetStatements([]string{"FOO"})
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "set -e ") {
		t.Errorf("unexpected line: %s", lines[0])
	}
}

func TestEvalBlockEmpty(t *testing.T) {
	e := NewExporter("bash")
	result := e.EvalBlock([]string{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestEvalBlockJoinsLines(t *testing.T) {
	e := NewExporter("bash")
	stmts := []string{"export A=\"1\"", "export B=\"2\""}
	result := e.EvalBlock(stmts)
	if !strings.Contains(result, "\n") {
		t.Errorf("expected newline-joined block, got %q", result)
	}
}
