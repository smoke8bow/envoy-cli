package shell

import (
	"fmt"
	"strings"
)

// ShellType represents a supported shell.
type ShellType string

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Fish ShellType = "fish"
)

// Exporter generates shell-specific export statements for environment variables.
type Exporter struct {
	Shell ShellType
}

// NewExporter creates an Exporter for the given shell type string.
// Defaults to bash if the shell is unrecognized.
func NewExporter(shell string) *Exporter {
	switch ShellType(strings.ToLower(shell)) {
	case Zsh:
		return &Exporter{Shell: Zsh}
	case Fish:
		return &Exporter{Shell: Fish}
	default:
		return &Exporter{Shell: Bash}
	}
}

// ExportStatements returns shell-specific export lines for the given env map.
func (e *Exporter) ExportStatements(env map[string]string) []string {
	lines := make([]string, 0, len(env))
	for k, v := range env {
		switch e.Shell {
		case Fish:
			lines = append(lines, fmt.Sprintf("set -x %s %q", k, v))
		default:
			lines = append(lines, fmt.Sprintf("export %s=%q", k, v))
		}
	}
	return lines
}

// UnsetStatements returns shell-specific unset lines for the given keys.
func (e *Exporter) UnsetStatements(keys []string) []string {
	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		switch e.Shell {
		case Fish:
			lines = append(lines, fmt.Sprintf("set -e %s", k))
		default:
			lines = append(lines, fmt.Sprintf("unset %s", k))
		}
	}
	return lines
}

// EvalBlock wraps statements into an eval-friendly block comment header.
func (e *Exporter) EvalBlock(statements []string) string {
	if len(statements) == 0 {
		return ""
	}
	return strings.Join(statements, "\n")
}
