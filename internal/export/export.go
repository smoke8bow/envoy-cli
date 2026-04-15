package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Format represents the output format for exported profiles.
type Format string

const (
	FormatJSON   Format = "json"
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
)

// ExportedProfile holds profile data for export.
type ExportedProfile struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	ExportedAt time.Time        `json:"exported_at"`
}

// Exporter writes profile data to files or stdout.
type Exporter struct {
	format Format
}

// NewExporter creates an Exporter for the given format.
func NewExporter(format Format) (*Exporter, error) {
	switch format {
	case FormatJSON, FormatDotenv, FormatShell:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// Render serializes the profile to the chosen format.
func (e *Exporter) Render(name string, vars map[string]string) (string, error) {
	ep := ExportedProfile{Name: name, Vars: vars, ExportedAt: time.Now().UTC()}
	switch e.format {
	case FormatJSON:
		b, err := json.MarshalIndent(ep, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	case FormatDotenv:
		return renderDotenv(vars), nil
	case FormatShell:
		return renderShell(vars), nil
	}
	return "", fmt.Errorf("unknown format")
}

// WriteFile writes the rendered output to dest path.
func (e *Exporter) WriteFile(dest, name string, vars map[string]string) error {
	out, err := e.Render(name, vars)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dest, []byte(out), 0o600)
}

func renderDotenv(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func renderShell(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}
