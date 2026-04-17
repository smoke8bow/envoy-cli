package namespace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/namespace"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "ns-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func newManager(t *testing.T) *namespace.FSManager {
	t.Helper()
	fs := namespace.NewFileStore(filepath.Join(tempDir(t)))
	return namespace.NewFSManager(fs)
}

func TestCreateAndList(t *testing.T) {
	m := newManager(t)
	if err := m.Create("staging"); err != nil {
		t.Fatal(err)
	}
	ns, err := m.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ns) != 1 || ns[0].Name != "staging" {
		t.Fatalf("expected staging, got %+v", ns)
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := newManager(t)
	m.Create("prod")
	if err := m.Create("prod"); err == nil {
		t.Fatal("expected error for duplicate")
	}
}

func TestDelete(t *testing.T) {
	m := newManager(t)
	m.Create("dev")
	if err := m.Delete("dev"); err != nil {
		t.Fatal(err)
	}
	ns, _ := m.List()
	if len(ns) != 0 {
		t.Fatal("expected empty list after delete")
	}
}

func TestDeleteNotFound(t *testing.T) {
	m := newManager(t)
	if err := m.Delete("ghost"); err == nil {
		t.Fatal("expected error")
	}
}

func TestAssignAndUnassign(t *testing.T) {
	m := newManager(t)
	m.Create("team")
	if err := m.Assign("team", "profile-a"); err != nil {
		t.Fatal(err)
	}
	ns, _ := m.List()
	if len(ns[0].Profiles) != 1 || ns[0].Profiles[0] != "profile-a" {
		t.Fatalf("unexpected profiles: %+v", ns[0].Profiles)
	}
	if err := m.Unassign("team", "profile-a"); err != nil {
		t.Fatal(err)
	}
	ns, _ = m.List()
	if len(ns[0].Profiles) != 0 {
		t.Fatal("expected empty profiles after unassign")
	}
}

func TestAssignIdempotent(t *testing.T) {
	m := newManager(t)
	m.Create("ci")
	m.Assign("ci", "base")
	m.Assign("ci", "base")
	ns, _ := m.List()
	if len(ns[0].Profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(ns[0].Profiles))
	}
}
