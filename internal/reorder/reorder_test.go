package reorder_test

import (
	"errors"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/reorder"
)

type fakeStore struct {
	profiles map[string]map[string]string
	saveErr  error
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	p, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := make(map[string]string, len(p))
	for k, v := range p {
		copy[k] = v
	}
	return copy, nil
}

func (f *fakeStore) Save(name string, vars map[string]string) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.profiles[name] = vars
	return nil
}

func newReorderer(vars map[string]string) (*reorder.Reorderer, *fakeStore) {
	fs := &fakeStore{
		profiles: map[string]map[string]string{"p": vars},
	}
	return reorder.NewReorderer(fs), fs
}

func TestApplyOrdering(t *testing.T) {
	r, _ := newReorderer(map[string]string{"C": "3", "A": "1", "B": "2"})
	result, err := r.Apply("p", []string{"A", "B", "C"})
	if err != nil {
		t.Fatal(err)
	}
	if result["A"] != "1" || result["B"] != "2" || result["C"] != "3" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestApplyUnknownKeysIgnored(t *testing.T) {
	r, _ := newReorderer(map[string]string{"A": "1"})
	result, err := r.Apply("p", []string{"Z", "A"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["Z"]; ok {
		t.Fatal("Z should not be in result")
	}
}

func TestApplyRemainingKeysAppended(t *testing.T) {
	r, _ := newReorderer(map[string]string{"A": "1", "B": "2", "C": "3"})
	result, err := r.Apply("p", []string{"A"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
}

func TestApplyProfileNotFound(t *testing.T) {
	r, _ := newReorderer(map[string]string{})
	_, err := r.Apply("missing", []string{"A"})
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestPreviewDoesNotPersist(t *testing.T) {
	original := map[string]string{"B": "2", "A": "1"}
	r, fs := newReorderer(original)
	keys, err := r.Preview("p", []string{"A", "B"})
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	// store should be unchanged
	if len(fs.profiles["p"]) != 2 {
		t.Fatal("store should not have been mutated")
	}
}

func TestApplySaveError(t *testing.T) {
	r, fs := newReorderer(map[string]string{"A": "1"})
	fs.saveErr = errors.New("disk full")
	_, err := r.Apply("p", []string{"A"})
	if err == nil {
		t.Fatal("expected save error")
	}
}
