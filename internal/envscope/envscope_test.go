package envscope

import (
	"testing"
)

func newManager(data map[string]map[string]string) *Manager {
	return NewManager(NewStoreAccessor(data))
}

func TestBuildEmptyScope(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"prod": {"APP_HOST": "localhost", "APP_PORT": "8080"},
	})
	_, err := m.Build("prod", "", false)
	if err == nil {
		t.Fatal("expected error for empty scope")
	}
}

func TestBuildEmptyProfile(t *testing.T) {
	m := newManager(map[string]map[string]string{})
	_, err := m.Build("", "APP_", false)
	if err == nil {
		t.Fatal("expected error for empty profile")
	}
}

func TestBuildProfileNotFound(t *testing.T) {
	m := newManager(map[string]map[string]string{})
	_, err := m.Build("missing", "APP_", false)
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestBuildNoStrip(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"prod": {
			"APP_HOST": "localhost",
			"APP_PORT": "8080",
			"DB_HOST":  "db",
		},
	})
	sv, err := m.Build("prod", "APP_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sv.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(sv.Vars))
	}
	if sv.Vars["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", sv.Vars["APP_HOST"])
	}
	if sv.Vars["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q", sv.Vars["APP_PORT"])
	}
}

func TestBuildWithStrip(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"prod": {
			"APP_HOST": "localhost",
			"APP_PORT": "8080",
			"DB_HOST":  "db",
		},
	})
	sv, err := m.Build("prod", "APP_", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sv.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(sv.Vars))
	}
	if sv.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", sv.Vars["HOST"])
	}
	if sv.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", sv.Vars["PORT"])
	}
}

func TestBuildCaseInsensitivePrefix(t *testing.T) {
	m := newManager(map[string]map[string]string{
		"dev": {
			"APP_HOST": "localhost",
			"OTHER":    "val",
		},
	})
	sv, err := m.Build("dev", "app_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sv.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(sv.Vars))
	}
}

func TestScopedViewKeys(t *testing.T) {
	sv := &ScopedView{
		Vars: map[string]string{"Z": "1", "A": "2", "M": "3"},
	}
	keys := sv.Keys()
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", keys)
	}
}
