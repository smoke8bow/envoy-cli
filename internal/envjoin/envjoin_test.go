package envjoin

import (
	"strings"
	"testing"
)

func TestJoinNoOverlap(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	out, err := Join(DefaultOptions(), a, b)
	if err != nil {
		t.Fatal(err)
	}
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Fatalf("unexpected result: %v", out)
	}
}

func TestJoinOverlappingKeys(t *testing.T) {
	a := map[string]string{"PATH": "/usr/bin"}
	b := map[string]string{"PATH": "/usr/local/bin"}
	out, err := Join(DefaultOptions(), a, b)
	if err != nil {
		t.Fatal(err)
	}
	want := "/usr/bin:/usr/local/bin"
	if out["PATH"] != want {
		t.Fatalf("got %q, want %q", out["PATH"], want)
	}
}

func TestJoinCustomSeparator(t *testing.T) {
	opts := Options{Separator: ",", Deduplicate: false}
	a := map[string]string{"TAGS": "alpha"}
	b := map[string]string{"TAGS": "beta"}
	out, err := Join(opts, a, b)
	if err != nil {
		t.Fatal(err)
	}
	if out["TAGS"] != "alpha,beta" {
		t.Fatalf("unexpected: %q", out["TAGS"])
	}
}

func TestJoinDeduplicate(t *testing.T) {
	opts := Options{Separator: ":", Deduplicate: true}
	a := map[string]string{"PATH": "/usr/bin"}
	b := map[string]string{"PATH": "/usr/bin"}
	c := map[string]string{"PATH": "/usr/local/bin"}
	out, err := Join(opts, a, b, c)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(out["PATH"], ":")
	if len(parts) != 2 {
		t.Fatalf("expected 2 unique parts, got %v", parts)
	}
}

func TestJoinEmptySeparatorError(t *testing.T) {
	opts := Options{Separator: ""}
	_, err := Join(opts, map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestJoinNoSources(t *testing.T) {
	out, err := Join(DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestFormat(t *testing.T) {
	sources := []map[string]string{
		{"A": "1", "PATH": "/usr/bin"},
		{"PATH": "/usr/local/bin"},
	}
	out, _ := Join(DefaultOptions(), sources...)
	s := Format(out, sources)
	if !strings.Contains(s, "2 total keys") {
		t.Fatalf("unexpected format output: %q", s)
	}
	if !strings.Contains(s, "1 joined") {
		t.Fatalf("expected joined count in output: %q", s)
	}
}
