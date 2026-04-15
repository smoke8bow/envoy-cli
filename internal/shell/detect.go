package shell

import (
	"os"
	"path/filepath"
	"strings"
)

// Detect attempts to determine the current shell from environment variables.
// It inspects $SHELL and falls back to "bash" if detection fails.
func Detect() string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return string(Bash)
	}

	base := strings.ToLower(filepath.Base(shellEnv))

	switch {
	case strings.Contains(base, "zsh"):
		return string(Zsh)
	case strings.Contains(base, "fish"):
		return string(Fish)
	case strings.Contains(base, "bash"):
		return string(Bash)
	default:
		return string(Bash)
	}
}

// Supported returns the list of shell names supported by the exporter.
func Supported() []string {
	return []string{
		string(Bash),
		string(Zsh),
		string(Fish),
	}
}

// IsSupported reports whether the given shell name is supported.
func IsSupported(shell string) bool {
	for _, s := range Supported() {
		if strings.EqualFold(s, shell) {
			return true
		}
	}
	return false
}
