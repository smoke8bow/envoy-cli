package protect_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/protect"
)

func newManager() *protect.Manager {
	return protect.NewManager()
}

func TestProtectAndIsProtected(t *testing.T) {
	m := newManager()
	if err := m.Protect("prod", "DB_PASSWORD"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.IsProtected("prod", "DB_PASSWORD") {
		t.Error("expected key to be protected")
	}
	if m.IsProtected("prod", "OTHER_KEY") {
		t.Error("expected OTHER_KEY to not be protected")
	}
}

func TestProtectEmptyProfileError(t *testing.T) {
	m := newManager()
	if err := m.Protect("", "KEY"); err == nil {
		t.Error("expected error for empty profile")
	}
}

func TestProtectEmptyKeyError(t *testing.T) {
	m := newManager()
	if err := m.Protect("prod", ""); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestUnprotect(t *testing.T) {
	m := newManager()
	_ = m.Protect("prod", "SECRET")
	if err := m.Unprotect("prod", "SECRET"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.IsProtected("prod", "SECRET") {
		t.Error("expected key to no longer be protected")
	}
}

func TestUnprotectNotProtected(t *testing.T) {
	m := newManager()
	err := m.Unprotect("prod", "MISSING")
	if !errors.Is(err, protect.ErrKeyNotProtected) {
		t.Errorf("expected ErrKeyNotProtected, got %v", err)
	}
}

func TestList(t *testing.T) {
	m := newManager()
	_ = m.Protect("prod", "A")
	_ = m.Protect("prod", "B")
	_ = m.Protect("prod", "C")
	keys := m.List("prod")
	sort.Strings(keys)
	if len(keys) != 3 || keys[0] != "A" || keys[1] != "B" || keys[2] != "C" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestListEmpty(t *testing.T) {
	m := newManager()
	if keys := m.List("nonexistent"); len(keys) != 0 {
		t.Errorf("expected empty list, got %v", keys)
	}
}

func TestGuardBlocks(t *testing.T) {
	m := newManager()
	_ = m.Protect("prod", "DB_PASS")
	err := m.Guard("prod", []string{"API_KEY", "DB_PASS"})
	if !errors.Is(err, protect.ErrKeyProtected) {
		t.Errorf("expected ErrKeyProtected, got %v", err)
	}
}

func TestGuardAllows(t *testing.T) {
	m := newManager()
	_ = m.Protect("prod", "DB_PASS")
	if err := m.Guard("prod", []string{"API_KEY", "HOST"}); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGuardNoProfile(t *testing.T) {
	m := newManager()
	if err := m.Guard("unknown", []string{"ANY_KEY"}); err != nil {
		t.Errorf("expected no error for unknown profile, got %v", err)
	}
}
