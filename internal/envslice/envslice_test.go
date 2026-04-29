package envslice_test

import (
	"testing"

	"envoy-cli/internal/envslice"
)

func TestFromMapSorted(t *testing.T) {
	vars := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	entries := envslice.FromMap(vars)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "APPLE" || entries[1].Key != "MANGO" || entries[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", entries)
	}
}

func TestFromMapEmpty(t *testing.T) {
	entries := envslice.FromMap(map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestToMapRoundtrip(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := envslice.ToMap(envslice.FromMap(vars))
	for k, v := range vars {
		if result[k] != v {
			t.Errorf("key %s: expected %q got %q", k, v, result[k])
		}
	}
}

func TestToStrings(t *testing.T) {
	entries := []envslice.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	strs := envslice.ToStrings(entries)
	if strs[0] != "A=1" || strs[1] != "B=2" {
		t.Errorf("unexpected strings: %v", strs)
	}
}

func TestFromStringsBasic(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux=extra"}
	entries := envslice.FromStrings(lines)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry[0]: %+v", entries[0])
	}
	if entries[1].Key != "BAZ" || entries[1].Value != "qux=extra" {
		t.Errorf("unexpected entry[1]: %+v", entries[1])
	}
}

func TestFromStringsSkipsNoEquals(t *testing.T) {
	lines := []string{"NOEQUALS", "KEY=val"}
	entries := envslice.FromStrings(lines)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "KEY" {
		t.Errorf("unexpected key: %s", entries[0].Key)
	}
}

func TestFilterByPrefix(t *testing.T) {
	entries := []envslice.Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db"},
	}
	result := envslice.FilterByPrefix(entries, "APP_")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	for _, e := range result {
		if e.Key[:4] != "APP_" {
			t.Errorf("unexpected key: %s", e.Key)
		}
	}
}

func TestFilterByPrefixNoMatch(t *testing.T) {
	entries := []envslice.Entry{{Key: "FOO", Value: "bar"}}
	result := envslice.FilterByPrefix(entries, "MISSING_")
	if len(result) != 0 {
		t.Errorf("expected empty, got %d entries", len(result))
	}
}

func TestEntryString(t *testing.T) {
	e := envslice.Entry{Key: "MY_VAR", Value: "hello world"}
	if e.String() != "MY_VAR=hello world" {
		t.Errorf("unexpected string: %s", e.String())
	}
}
