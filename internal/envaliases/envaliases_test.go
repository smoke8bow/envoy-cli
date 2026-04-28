package envaliases_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envaliases"
)

func newManager() *envaliases.Manager {
	return envaliases.NewManager()
}

func TestSetAndResolve(t *testing.T) {
	m := newManager()
	if err := m.Set("prod", "db", "DATABASE_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	key, err := m.Resolve("prod", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %q", key)
	}
}

func TestSetEmptyProfileError(t *testing.T) {
	m := newManager()
	if err := m.Set("", "db", "DATABASE_URL"); err == nil {
		t.Error("expected error for empty profile")
	}
}

func TestSetEmptyAliasError(t *testing.T) {
	m := newManager()
	if err := m.Set("prod", "", "DATABASE_URL"); err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	m := newManager()
	if err := m.Set("prod", "db", ""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestResolveUnknownProfile(t *testing.T) {
	m := newManager()
	_, err := m.Resolve("ghost", "db")
	if err == nil {
		t.Error("expected error for unknown profile")
	}
}

func TestResolveUnknownAlias(t *testing.T) {
	m := newManager()
	_ = m.Set("prod", "db", "DATABASE_URL")
	_, err := m.Resolve("prod", "cache")
	if err == nil {
		t.Error("expected error for unknown alias")
	}
}

func TestList(t *testing.T) {
	m := newManager()
	_ = m.Set("prod", "db", "DATABASE_URL")
	_ = m.Set("prod", "cache", "REDIS_URL")
	list := m.List("prod")
	if len(list) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(list))
	}
	if list[0] != "cache" || list[1] != "db" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestRemove(t *testing.T) {
	m := newManager()
	_ = m.Set("prod", "db", "DATABASE_URL")
	if err := m.Remove("prod", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.List("prod")) != 0 {
		t.Error("expected empty list after remove")
	}
}

func TestRemoveNotFound(t *testing.T) {
	m := newManager()
	if err := m.Remove("prod", "db"); err == nil {
		t.Error("expected error removing from unknown profile")
	}
}

func TestExpand(t *testing.T) {
	m := newManager()
	_ = m.Set("prod", "db", "DATABASE_URL")
	_ = m.Set("prod", "cache", "REDIS_URL")
	input := map[string]string{
		"db":    "postgres://localhost/mydb",
		"cache": "redis://localhost:6379",
		"OTHER": "value",
	}
	out := m.Expand("prod", input)
	if out["DATABASE_URL"] != "postgres://localhost/mydb" {
		t.Errorf("expected DATABASE_URL to be expanded")
	}
	if out["REDIS_URL"] != "redis://localhost:6379" {
		t.Errorf("expected REDIS_URL to be expanded")
	}
	if out["OTHER"] != "value" {
		t.Errorf("expected OTHER to pass through")
	}
}

func TestExpandDoesNotMutateInput(t *testing.T) {
	m := newManager()
	_ = m.Set("prod", "db", "DATABASE_URL")
	input := map[string]string{"db": "postgres://localhost"}
	_ = m.Expand("prod", input)
	if _, ok := input["DATABASE_URL"]; ok {
		t.Error("input map was mutated")
	}
}
