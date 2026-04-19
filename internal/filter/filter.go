package filter

import (
	"strings"
)

// Option configures filtering behaviour.
type Option struct {
	Prefix    string
	Suffix    string
	Contains  string
	ExactKeys []string
}

// Result holds the filtered key/value pairs and the keys that were excluded.
type Result struct {
	Matched  map[string]string
	Excluded []string
}

// Filter applies the given Option to vars and returns a Result.
func Filter(vars map[string]string, opt Option) Result {
	matched := make(map[string]string)
	var excluded []string

	exactSet := make(map[string]struct{}, len(opt.ExactKeys))
	for _, k := range opt.ExactKeys {
		exactSet[k] = struct{}{}
	}

	for k, v := range vars {
		if matches(k, opt, exactSet) {
			matched[k] = v
		} else {
			excluded = append(excluded, k)
		}
	}
	return Result{Matched: matched, Excluded: excluded}
}

func matches(key string, opt Option, exactSet map[string]struct{}) bool {
	if len(exactSet) > 0 {
		_, ok := exactSet[key]
		return ok
	}
	if opt.Prefix != "" && !strings.HasPrefix(key, opt.Prefix) {
		return false
	}
	if opt.Suffix != "" && !strings.HasSuffix(key, opt.Suffix) {
		return false
	}
	if opt.Contains != "" && !strings.Contains(key, opt.Contains) {
		return false
	}
	return true
}
