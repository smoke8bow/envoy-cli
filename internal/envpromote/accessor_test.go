package envpromote

import (
	"errors"
	"testing"
)

func TestStoreAccessorGet(t *testing.T) {
	want := map[string]string{"KEY": "val"}
	acc := NewStoreAccessor(
		func(name string) (map[string]string, error) {
			if name == "myprofile" {
				return want, nil
			}
			return nil, errors.New("not found")
		},
		func(string, map[string]string) error { return nil },
	)

	got, err := acc.Get("myprofile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", got["KEY"])
	}
}

func TestStoreAccessorGetError(t *testing.T) {
	acc := NewStoreAccessor(
		func(string) (map[string]string, error) { return nil, errors.New("boom") },
		func(string, map[string]string) error { return nil },
	)
	_, err := acc.Get("x")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStoreAccessorSet(t *testing.T) {
	var captured map[string]string
	acc := NewStoreAccessor(
		func(string) (map[string]string, error) { return nil, nil },
		func(_ string, vars map[string]string) error {
			captured = vars
			return nil
		},
	)
	input := map[string]string{"A": "1"}
	if err := acc.Set("p", input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured["A"] != "1" {
		t.Errorf("expected A=1 in captured vars")
	}
}

func TestStoreAccessorSetError(t *testing.T) {
	acc := NewStoreAccessor(
		func(string) (map[string]string, error) { return nil, nil },
		func(string, map[string]string) error { return errors.New("write fail") },
	)
	if err := acc.Set("p", nil); err == nil {
		t.Fatal("expected error")
	}
}
