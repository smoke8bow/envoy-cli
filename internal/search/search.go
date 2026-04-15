package search

import (
	"sort"
	"strings"
)

// Result holds a matched profile name and the keys that matched the query.
type Result struct {
	Profile string
	MatchedKeys []string
}

// Source provides profile data for searching.
type Source interface {
	List() ([]string, error)
	Get(name string) (map[string]string, error)
}

// Searcher performs queries against profile environment variables.
type Searcher struct {
	src Source
}

// NewSearcher creates a Searcher backed by the given Source.
func NewSearcher(src Source) *Searcher {
	return &Searcher{src: src}
}

// ByKey returns all profiles that contain a key matching the given substring
// (case-insensitive). Results are sorted by profile name.
func (s *Searcher) ByKey(query string) ([]Result, error) {
	names, err := s.src.List()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var results []Result

	for _, name := range names {
		vars, err := s.src.Get(name)
		if err != nil {
			return nil, err
		}
		var matched []string
		for k := range vars {
			if strings.Contains(strings.ToLower(k), q) {
				matched = append(matched, k)
			}
		}
		if len(matched) > 0 {
			sort.Strings(matched)
			results = append(results, Result{Profile: name, MatchedKeys: matched})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Profile < results[j].Profile
	})
	return results, nil
}

// ByValue returns all profiles that contain a value matching the given
// substring (case-insensitive). Results are sorted by profile name.
func (s *Searcher) ByValue(query string) ([]Result, error) {
	names, err := s.src.List()
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var results []Result

	for _, name := range names {
		vars, err := s.src.Get(name)
		if err != nil {
			return nil, err
		}
		var matched []string
		for k, v := range vars {
			if strings.Contains(strings.ToLower(v), q) {
				matched = append(matched, k)
			}
		}
		if len(matched) > 0 {
			sort.Strings(matched)
			results = append(results, Result{Profile: name, MatchedKeys: matched})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Profile < results[j].Profile
	})
	return results, nil
}
