package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Resolver is a function that returns the value for a given key.
type Resolver func(key string) (string, bool)

// MapResolver returns a Resolver backed by a static map.
func MapResolver(m map[string]string) Resolver {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

// EnvResolver returns a Resolver that falls back to the process environment.
func EnvResolver() Resolver {
	return func(key string) (string, bool) {
		v, ok := os.LookupEnv(key)
		return v, ok
	}
}

// ChainResolver tries each resolver in order, returning the first match.
func ChainResolver(resolvers ...Resolver) Resolver {
	return func(key string) (string, bool) {
		for _, r := range resolvers {
			if v, ok := r(key); ok {
				return v, true
			}
		}
		return "", false
	}
}

// Expand replaces variable references in s using the provided Resolver.
// Unknown variables are left as-is.
func Expand(s string, resolve Resolver) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		key := extractKey(match)
		if v, ok := resolve(key); ok {
			return v
		}
		return match
	})
}

// ExpandStrict is like Expand but returns an error if any variable is unresolved.
func ExpandStrict(s string, resolve Resolver) (string, error) {
	var missing []string
	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		key := extractKey(match)
		if v, ok := resolve(key); ok {
			return v
		}
		missing = append(missing, key)
		return match
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("unresolved variables: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// ExpandMap applies Expand to every value in the map, returning a new map.
func ExpandMap(vars map[string]string, resolve Resolver) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = Expand(v, resolve)
	}
	return out
}

func extractKey(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
