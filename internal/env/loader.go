package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadDotenv reads a .env file from the given path and returns a map of
// key=value pairs. Lines beginning with '#' are treated as comments and
// skipped. Blank lines are also skipped. Values may optionally be quoted
// with single or double quotes; the quotes are stripped before returning.
func LoadDotenv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("env: open %q: %w", path, err)
	}
	defer f.Close()

	vars := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip inline export keyword (e.g. "export FOO=bar").
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("env: %q line %d: missing '=' in %q", path, lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		if key == "" {
			return nil, fmt.Errorf("env: %q line %d: empty key", path, lineNum)
		}

		val = stripQuotes(val)
		vars[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scan %q: %w", path, err)
	}

	return vars, nil
}

// stripQuotes removes surrounding single or double quotes from s, if present.
// It only strips when both the opening and closing quote match.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// LoadOS reads all current process environment variables and returns them
// as a map. This is a convenience wrapper around os.Environ.
func LoadOS() map[string]string {
	result := make(map[string]string)
	for _, entry := range os.Environ() {
		idx := strings.IndexByte(entry, '=')
		if idx < 0 {
			continue
		}
		result[entry[:idx]] = entry[idx+1:]
	}
	return result
}
