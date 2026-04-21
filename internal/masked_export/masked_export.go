// Package masked_export provides export functionality that automatically
// redacts sensitive environment variable values before rendering output.
package masked_export

import (
	"fmt"

	"github.com/your-org/envoy-cli/internal/mask"
)

// Format represents a supported export format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// Exporter renders masked environment variable maps.
type Exporter struct {
	masker  *mask.Masker
	format  Format
	reveal  []string // keys whose values should NOT be masked
}

// NewExporter creates a new masked exporter for the given format.
// An error is returned if the format is unsupported.
func NewExporter(format Format, patterns []string, reveal []string) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatShell, FormatJSON:
		// valid
	default:
		return nil, fmt.Errorf("unsupported format %q: must be dotenv, shell, or json", format)
	}

	var m *mask.Masker
	if len(patterns) > 0 {
		m = mask.NewMasker(patterns)
	} else {
		m = mask.NewMasker(nil)
	}

	return &Exporter{
		masker: m,
		format: format,
		reveal: reveal,
	}, nil
}

// Render applies masking to vars and returns the formatted output string.
func (e *Exporter) Render(vars map[string]string) (string, error) {
	masked := e.masker.Apply(vars)
	for _, key := range e.reveal {
		if orig, ok := vars[key]; ok {
			masked[key] = orig
		}
	}

	switch e.format {
	case FormatDotenv:
		return renderDotenv(masked), nil
	case FormatShell:
		return renderShell(masked), nil
	case FormatJSON:
		return renderJSON(masked), nil
	}
	return "", fmt.Errorf("unsupported format")
}

func renderDotenv(vars map[string]string) string {
	out := ""
	for k, v := range vars {
		out += fmt.Sprintf("%s=%q\n", k, v)
	}
	return out
}

func renderShell(vars map[string]string) string {
	out := ""
	for k, v := range vars {
		out += fmt.Sprintf("export %s=%q\n", k, v)
	}
	return out
}

func renderJSON(vars map[string]string) string {
	out := "{"
	first := true
	for k, v := range vars {
		if !first {
			out += ","
		}
		out += fmt.Sprintf("%q:%q", k, v)
		first = false
	}
	out += "}"
	return out
}
