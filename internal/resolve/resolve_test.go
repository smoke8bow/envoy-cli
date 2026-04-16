package resolve

import (
	"testing"
)

func TestResolveNoRefs(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestResolveBraceRef(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"BASE": "/usr", "BIN": "${BASE}/bin"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["BIN"] != "/usr/bin" {
		t.Errorf("got %q", out["BIN"])
	}
}

func TestResolveDollarRef(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"HOST": "localhost", "ADDR": "$HOST:8080"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["ADDR"] != "localhost:8080" {
		t.Errorf("got %q", out["ADDR"])
	}
}

func TestResolveAmbient(t *testing.T) {
	ambient := map[string]string{"HOME": "/home/user"}
	r := NewResolver(ambient)
	in := map[string]string{"CONF": "${HOME}/.config"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["CONF"] != "/home/user/.config" {
		t.Errorf("got %q", out["CONF"])
	}
}

func TestResolveUndefinedError(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"X": "${MISSING}"}
	_, err := r.Resolve(in)
	if err == nil {
		t.Fatal("expected error for undefined variable")
	}
}

func TestResolveCycleError(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"A": "${B}", "B": "${A}"}
	_, err := r.Resolve(in)
	if err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestUnresolvedKeys(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"GOOD": "plain", "BAD": "${NOPE}"}
	keys := r.UnresolvedKeys(in)
	if len(keys) != 1 || keys[0] != "BAD" {
		t.Errorf("unexpected unresolved keys: %v", keys)
	}
}

func TestResolveChained(t *testing.T) {
	r := NewResolver(nil)
	in := map[string]string{"A": "hello", "B": "${A}_world", "C": "${B}!"}
	out, err := r.Resolve(in)
	if err != nil {
		t.Fatal(err)
	}
	if out["C"] != "hello_world!" {
		t.Errorf("got %q", out["C"])
	}
}
