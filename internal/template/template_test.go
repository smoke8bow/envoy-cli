package template_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/template"
)

func TestExpandBraceStyle(t *testing.T) {
	resolve := template.MapResolver(map[string]string{"HOME": "/home/user"})
	got := template.Expand("path: ${HOME}/bin", resolve)
	if got != "path: /home/user/bin" {
		t.Fatalf("expected 'path: /home/user/bin', got %q", got)
	}
}

func TestExpandDollarStyle(t *testing.T) {
	resolve := template.MapResolver(map[string]string{"USER": "alice"})
	got := template.Expand("hello $USER", resolve)
	if got != "hello alice" {
		t.Fatalf("expected 'hello alice', got %q", got)
	}
}

func TestExpandUnknownLeft(t *testing.T) {
	resolve := template.MapResolver(map[string]string{})
	got := template.Expand("${MISSING}", resolve)
	if got != "${MISSING}" {
		t.Fatalf("expected original token, got %q", got)
	}
}

func TestExpandStrictSuccess(t *testing.T) {
	resolve := template.MapResolver(map[string]string{"PORT": "8080"})
	got, err := template.ExpandStrict("port=${PORT}", resolve)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "port=8080" {
		t.Fatalf("expected 'port=8080', got %q", got)
	}
}

func TestExpandStrictMissing(t *testing.T) {
	resolve := template.MapResolver(map[string]string{})
	_, err := template.ExpandStrict("${MISSING}", resolve)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExpandMap(t *testing.T) {
	base := map[string]string{
		"DB_URL": "postgres://${DB_HOST}:${DB_PORT}/mydb",
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	resolve := template.MapResolver(base)
	out := template.ExpandMap(base, resolve)
	want := "postgres://localhost:5432/mydb"
	if out["DB_URL"] != want {
		t.Fatalf("expected %q, got %q", want, out["DB_URL"])
	}
}

func TestChainResolver(t *testing.T) {
	primary := template.MapResolver(map[string]string{"A": "from-primary"})
	fallback := template.MapResolver(map[string]string{"A": "from-fallback", "B": "from-fallback"})
	chain := template.ChainResolver(primary, fallback)

	if v, ok := chain("A"); !ok || v != "from-primary" {
		t.Fatalf("expected primary value for A, got %q", v)
	}
	if v, ok := chain("B"); !ok || v != "from-fallback" {
		t.Fatalf("expected fallback value for B, got %q", v)
	}
	if _, ok := chain("C"); ok {
		t.Fatal("expected miss for C")
	}
}

func TestExpandMultipleVars(t *testing.T) {
	resolve := template.MapResolver(map[string]string{
		"PROTO": "https",
		"HOST":  "example.com",
	})
	got := template.Expand("${PROTO}://${HOST}/api", resolve)
	want := "https://example.com/api"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
