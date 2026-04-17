package expire

import (
	"os"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "expire-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newManager(t *testing.T) *Manager {
	t.Helper()
	m, err := NewManager(tempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestCheckNoExpiry(t *testing.T) {
	m := newManager(t)
	if err := m.Check("dev"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestSetAndCheckNotExpired(t *testing.T) {
	m := newManager(t)
	if err := m.Set("dev", 10*time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := m.Check("dev"); err != nil {
		t.Fatalf("expected not expired, got %v", err)
	}
}

func TestSetAndCheckExpired(t *testing.T) {
	m := newManager(t)
	if err := m.Set("dev", -1*time.Second); err != nil {
		t.Fatal(err)
	}
	if err := m.Check("dev"); err != ErrExpired {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestClearRemovesExpiry(t *testing.T) {
	m := newManager(t)
	m.Set("dev", -1*time.Second)
	if err := m.Clear("dev"); err != nil {
		t.Fatal(err)
	}
	if err := m.Check("dev"); err != nil {
		t.Fatalf("expected nil after clear, got %v", err)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	m1, _ := NewManager(dir)
	m1.Set("prod", 10*time.Minute)

	m2, err := NewManager(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := m2.Check("prod"); err != nil {
		t.Fatalf("expected valid after reload, got %v", err)
	}
}
