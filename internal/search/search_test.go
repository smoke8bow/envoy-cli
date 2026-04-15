package search

import (
	"errors"
	"testing"
)

// fakeSource implements Source for testing.
type fakeSource struct {
	profiles map[string]map[string]string
	listErr  error
	getErr   error
}

func (f *fakeSource) List() ([]string, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	names := make([]string, 0, len(f.profiles))
	for k := range f.profiles {
		names = append(names, k)
	}
	return names, nil
}

func (f *fakeSource) Get(name string) (map[string]string, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.profiles[name], nil
}

func newSearcher() *Searcher {
	src := &fakeSource{
		profiles: map[string]map[string]string{
			"dev": {"DB_HOST": "localhost", "API_KEY": "dev-secret"},
			"prod": {"DB_HOST": "prod.db", "LOG_LEVEL": "warn"},
			"staging": {"API_KEY": "stg-secret", "FEATURE_FLAG": "true"},
		},
	}
	return NewSearcher(src)
}

func TestByKeyMatch(t *testing.T) {
	s := newSearcher()
	results, err := s.ByKey("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Profile != "dev" || results[1].Profile != "staging" {
		t.Errorf("unexpected profiles: %v", results)
	}
}

func TestByKeyNoMatch(t *testing.T) {
	s := newSearcher()
	results, err := s.ByKey("NOPE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestByValueMatch(t *testing.T) {
	s := newSearcher()
	results, err := s.ByValue("secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestByValueCaseInsensitive(t *testing.T) {
	s := newSearcher()
	results, err := s.ByValue("WARN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Profile != "prod" {
		t.Errorf("expected prod, got %v", results)
	}
}

func TestByKeyListError(t *testing.T) {
	src := &fakeSource{listErr: errors.New("store unavailable")}
	s := NewSearcher(src)
	_, err := s.ByKey("anything")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestByValueGetError(t *testing.T) {
	src := &fakeSource{
		profiles: map[string]map[string]string{"dev": {"K": "v"}},
		getErr:   errors.New("read error"),
	}
	s := NewSearcher(src)
	_, err := s.ByValue("v")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
