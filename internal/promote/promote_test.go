package promote

import (
	"errors"
	"sort"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (f *fakeStore) Save(name string, vars map[string]string) error {
	f.profiles[name] = vars
	return nil
}

func (f *fakeStore) List() ([]string, error) {
	var names []string
	for k := range f.profiles {
		names = append(names, k)
	}
	return names, nil
}

func newPromoter() (*Promoter, *fakeStore) {
	fs := &fakeStore{profiles: map[string]map[string]string{
		"staging": {"DB_HOST": "stage-db", "API_KEY": "abc123", "DEBUG": "true"},
		"prod":    {"DB_HOST": "prod-db"},
	}}
	return NewPromoter(fs), fs
}

func TestPromoteAllKeys(t *testing.T) {
	p, fs := newPromoter()
	written, err := p.Promote("staging", "prod", PromoteOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 3 {
		t.Fatalf("expected 3 written, got %d", len(written))
	}
	if fs.profiles["prod"]["API_KEY"] != "abc123" {
		t.Error("API_KEY not promoted")
	}
}

func TestPromoteSelectedKeys(t *testing.T) {
	p, fs := newPromoter()
	written, err := p.Promote("staging", "prod", PromoteOptions{Keys: []string{"API_KEY"}, Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != 1 || written[0] != "API_KEY" {
		t.Errorf("unexpected written: %v", written)
	}
	if fs.profiles["prod"]["DB_HOST"] != "prod-db" {
		t.Error("DB_HOST should not have been overwritten")
	}
}

func TestPromoteNoOverwrite(t *testing.T) {
	p, _ := newPromoter()
	written, err := p.Promote("staging", "prod", PromoteOptions{Overwrite: false})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(written)
	// DB_HOST exists in prod and should be skipped
	for _, k := range written {
		if k == "DB_HOST" {
			t.Error("DB_HOST should have been skipped")
		}
	}
}

func TestPromoteSourceNotFound(t *testing.T) {
	p, _ := newPromoter()
	_, err := p.Promote("missing", "prod", PromoteOptions{})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestPromoteMissingKey(t *testing.T) {
	p, _ := newPromoter()
	_, err := p.Promote("staging", "prod", PromoteOptions{Keys: []string{"NONEXISTENT"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
