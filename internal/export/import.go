package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ImportFromFile reads a previously exported file and returns the profile name
// and variable map. Only JSON and dotenv formats are supported for import.
func ImportFromFile(path string) (name string, vars map[string]string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("reading file: %w", err)
	}
	if isJSON(data) {
		return parseJSON(data)
	}
	return "", parseDotenv(string(data)), nil
}

func isJSON(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	return len(trimmed) > 0 && trimmed[0] == '{'
}

func parseJSON(data []byte) (string, map[string]string, error) {
	var ep ExportedProfile
	if err := json.Unmarshal(data, &ep); err != nil {
		return "", nil, fmt.Errorf("parsing JSON: %w", err)
	}
	if ep.Vars == nil {
		ep.Vars = map[string]string{}
	}
	return ep.Name, ep.Vars, nil
}

func parseDotenv(content string) map[string]string {
	vars := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		vars[key] = val
	}
	return vars
}
