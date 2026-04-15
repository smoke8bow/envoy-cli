package env

import (
	"fmt"
	"os"
	"strings"
)

// Apply sets the given environment variables in the current process.
func Apply(vars map[string]string) error {
	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set env var %q: %w", key, err)
		}
	}
	return nil
}

// Export returns a shell-compatible export string for the given vars.
// Each line is in the form: export KEY="VALUE"
func Export(vars map[string]string) string {
	var sb strings.Builder
	for key, value := range vars {
		// Escape double quotes in value
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		sb.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, escaped))
	}
	return sb.String()
}

// Unset removes the given environment variable keys from the current process.
func Unset(keys []string) error {
	for _, key := range keys {
		if err := os.Unsetenv(key); err != nil {
			return fmt.Errorf("failed to unset env var %q: %w", key, err)
		}
	}
	return nil
}

// Snapshot captures the current values of the given keys from the environment.
// Keys not present in the environment are omitted from the result.
func Snapshot(keys []string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			result[key] = val
		}
	}
	return result
}
