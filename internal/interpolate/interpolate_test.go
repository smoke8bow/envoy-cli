package interpolate

import (
	"errors"
	"testing"
)

type fakeGetter struct {
	profiles map[string]map[string]string
}

func (f *fakeGetter) Get(name string) (map[string]string, error) {
	p, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("profile not found: " + name)
	}
	return p, nil
}

func newInterp(vars map[string]string) *Interpolator {
	g := &fakeGetter{profiles: map[string]map[string]string{"p": vars}}
	return New(g, DefaultOptions())
}

func TestNoReferences(t *testing.T) {
	interp := newInterp(map[string]string{"FOO": "bar"})
	out, err := interp.Apply(map[string]string{"FOO": "bar"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "bar" {
		t.Fatalf("expected bar, got %s", out["FOO"])
	}
}

func TestBraceStyleRef(t *testing.T) {
	vars := map[string]string{"BASE": "/usr", "BIN": "${BASE}/bin"}
	interp := newInterp(vars)
	out, err := interp.Apply(vars, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["BIN"] != "/usr/bin" {
		t.Fatalf("expected /usr/bin, got %s", out["BIN"])
	}
}

func TestDollarStyleRef(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "URL": "http://$HOST:8080"}
	interp := newInterp(vars)
	out, err := interp.Apply(vars, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Fatalf("expected http://localhost:8080, got %s", out["URL"])
	}
}

func TestAmbientFallback(t *testing.T) {
	vars := map[string]string{"GREETING": "hello ${EXTERNAL}"}
	interp := newInterp(vars)
	ambient := func(k string) string {
		if k == "EXTERNAL" {
			return "world"
		}
		return ""
	}
	out, err := interp.Apply(vars, ambient)
	if err != nil {
		t.Fatal(err)
	}
	if out["GREETING"] != "hello world" {
		t.Fatalf("expected 'hello world', got %s", out["GREETING"])
	}
}

func TestUnresolvedReturnsError(t *testing.T) {
	opts := DefaultOptions()
	opts.FallbackToOS = false
	vars := map[string]string{"A": "${MISSING}"}
	g := &fakeGetter{profiles: map[string]map[string]string{"p": vars}}
	interp := New(g, opts)
	_, err := interp.Apply(vars, nil)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestChainedRefs(t *testing.T) {
	vars := map[string]string{
		"A": "alpha",
		"B": "${A}-beta",
		"C": "${B}-gamma",
	}
	interp := newInterp(vars)
	out, err := interp.Apply(vars, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out["C"] != "alpha-beta-gamma" {
		t.Fatalf("expected alpha-beta-gamma, got %s", out["C"])
	}
}

func TestDoesNotMutateInput(t *testing.T) {
	vars := map[string]string{"X": "${Y}", "Y": "resolved"}
	orig := map[string]string{"X": "${Y}", "Y": "resolved"}
	interp := newInterp(vars)
	_, _ = interp.Apply(vars, nil)
	for k, v := range orig {
		if vars[k] != v {
			t.Fatalf("input mutated at key %s", k)
		}
	}
}
