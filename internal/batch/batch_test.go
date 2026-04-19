package batch

import (
	"errors"
	"testing"
)

type fakeStore struct {
	profiles map[string]map[string]string
	saveFail bool
}

func (f *fakeStore) Get(profile string) (map[string]string, error) {
	v, ok := f.profiles[profile]
	if !ok {
		return nil, errors.New("profile not found")
	}
	out := make(map[string]string, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out, nil
}

func (f *fakeStore) Save(profile string, vars map[string]string) error {
	if f.saveFail {
		return errors.New("save failed")
	}
	f.profiles[profile] = vars
	return nil
}

func newProcessor() (*Processor, *fakeStore) {
	s := &fakeStore{profiles: map[string]map[string]string{
		"dev": {"HOST": "localhost", "PORT": "8080"},
	}}
	return NewProcessor(s), s
}

func TestApplySet(t *testing.T) {
	p, s := newProcessor()
	ops := []Op{{Key: "DEBUG", Value: "true", Kind: OpSet}}
	_, err := p.Apply("dev", ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.profiles["dev"]["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true")
	}
}

func TestApplyDelete(t *testing.T) {
	p, s := newProcessor()
	ops := []Op{{Key: "PORT", Kind: OpDelete}}
	_, err := p.Apply("dev", ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.profiles["dev"]["PORT"]; ok {
		t.Errorf("expected PORT to be deleted")
	}
}

func TestApplyDeleteNotFound(t *testing.T) {
	p, _ := newProcessor()
	ops := []Op{{Key: "MISSING", Kind: OpDelete}}
	_, err := p.Apply("dev", ops)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestApplyEmptyKey(t *testing.T) {
	p, _ := newProcessor()
	ops := []Op{{Key: "", Value: "x", Kind: OpSet}}
	_, err := p.Apply("dev", ops)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestApplyProfileNotFound(t *testing.T) {
	p, _ := newProcessor()
	_, err := p.Apply("nope", []Op{{Key: "X", Value: "1", Kind: OpSet}})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestApplyDoesNotMutateOnError(t *testing.T) {
	p, s := newProcessor()
	ops := []Op{
		{Key: "NEW", Value: "val", Kind: OpSet},
		{Key: "MISSING", Kind: OpDelete},
	}
	p.Apply("dev", ops)
	if _, ok := s.profiles["dev"]["NEW"]; ok {
		t.Errorf("profile should not be mutated after error")
	}
}
