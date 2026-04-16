package resolve

import (
	"fmt"
	"sort"
	"strings"
)

// Resolver resolves environment variable references within a profile's values.
// References use the ${VAR} or $VAR syntax and can point to other keys in the
// same profile or to a base map of ambient variables.
type Resolver struct {
	ambient map[string]string
	maxDepth int
}

// NewResolver creates a Resolver with optional ambient variables.
func NewResolver(ambient map[string]string) *Resolver {
	if ambient == nil {
		ambient = map[string]string{}
	}
	return &Resolver{ambient: ambient, maxDepth: 10}
}

// Resolve expands all variable references in vars, using vars itself plus
// ambient variables as the source. Returns an error on cycles or missing keys.
func (r *Resolver) Resolve(vars map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		resolved, err := r.expand(k, v, vars, map[string]bool{}, 0)
		if err != nil {
			return nil, err
		}
		out[k] = resolved
	}
	return out, nil
}

// UnresolvedKeys returns keys whose values contain unresolvable references.
func (r *Resolver) UnresolvedKeys(vars map[string]string) []string {
	_, err := r.Resolve(vars)
	if err == nil {
		return nil
	}
	var bad []string
	for k, v := range vars {
		if strings.Contains(v, "$") {
			bad = append(bad, k)
		}
	}
	sort.Strings(bad)
	return bad
}

func (r *Resolver) expand(key, value string, vars map[string]string, visited map[string]bool, depth int) (string, error) {
	if depth > r.maxDepth {
		return "", fmt.Errorf("resolve: cycle or depth exceeded for key %q", key)
	}
	result := value
	for {
		idx := strings.Index(result, "$")
		if idx == -1 {
			break
		}
		name, length := extractVarName(result[idx:])
		if name == "" {
			break
		}
		val, ok := vars[name]
		if !ok {
			val, ok = r.ambient[name]
		}
		if !ok {
			return "", fmt.Errorf("resolve: undefined variable %q referenced in %q", name, key)
		}
		if visited[name] {
			return "", fmt.Errorf("resolve: cycle detected at %q", name)
		}
		visited[name] = true
		expanded, err := r.expand(name, val, vars, visited, depth+1)
		delete(visited, name)
		if err != nil {
			return "", err
		}
		result = result[:idx] + expanded + result[idx+length:]
	}
	return result, nil
}

func extractVarName(s string) (string, int) {
	if len(s) < 2 {
		return "", 0
	}
	if s[1] == '{' {
		end := strings.Index(s, "}")
		if end == -1 {
			return "", 0
		}
		return s[2:end], end + 1
	}
	i := 1
	for i < len(s) && (isAlnum(s[i]) || s[i] == '_') {
		i++
	}
	if i == 1 {
		return "", 0
	}
	return s[1:i], i
}

func isAlnum(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_'
}
