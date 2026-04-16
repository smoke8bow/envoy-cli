package inject

import (
	"strings"
	"testing"
)

func TestBuildOverlayOverridesBase(t *testing.T) {
	inj := NewInjector([]string{"FOO=base", "BAR=keep"})
	env := inj.Build(map[string]string{"FOO": "overridden"})
	m := toMap(env)
	if m["FOO"] != "overridden" {
		t.Fatalf("expected FOO=overridden, got %s", m["FOO"])
	}
	if m["BAR"] != "keep" {
		t.Fatalf("expected BAR=keep, got %s", m["BAR"])
	}
}

func TestBuildAddsNewKeys(t *testing.T) {
	inj := NewInjector([]string{"EXISTING=yes"})
	env := inj.Build(map[string]string{"NEW_KEY": "hello"})
	m := toMap(env)
	if m["NEW_KEY"] != "hello" {
		t.Fatalf("expected NEW_KEY=hello, got %s", m["NEW_KEY"])
	}
	if m["EXISTING"] != "yes" {
		t.Fatalf("expected EXISTING=yes, got %s", m["EXISTING"])
	}
}

func TestBuildEmptyOverlay(t *testing.T) {
	base := []string{"A=1", "B=2"}
	inj := NewInjector(base)
	env := inj.Build(nil)
	if len(env) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env))
	}
}

func TestBuildEmptyBase(t *testing.T) {
	inj := NewInjector(nil)
	env := inj.Build(map[string]string{"ONLY": "me"})
	m := toMap(env)
	if m["ONLY"] != "me" {
		t.Fatalf("expected ONLY=me, got %s", m["ONLY"])
	}
}

func TestCommandSetsEnv(t *testing.T) {
	inj := NewInjector([]string{"BASE=1"})
	cmd := inj.Command(map[string]string{"EXTRA": "2"}, "env")
	m := toMap(cmd.Env)
	if m["BASE"] != "1" {
		t.Fatalf("expected BASE=1")
	}
	if m["EXTRA"] != "2" {
		t.Fatalf("expected EXTRA=2")
	}
}

func toMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, kv := range env {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
