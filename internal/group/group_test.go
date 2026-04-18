package group_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/group"
)

type fakeStore struct{}

func (f *fakeStore) List() ([]string, error) { return []string{"dev", "prod"}, nil }

func newManager() *group.Manager {
	return group.NewManager(&fakeStore{})
}

func TestCreateAndList(t *testing.T) {
	m := newManager()
	if err := m.Create("team-a"); err != nil {
		t.Fatal(err)
	}
	if err := m.Create("team-b"); err != nil {
		t.Fatal(err)
	}
	names := m.List()
	if len(names) != 2 || names[0] != "team-a" || names[1] != "team-b" {
		t.Fatalf("unexpected list: %v", names)
	}
}

func TestCreateDuplicate(t *testing.T) {
	m := newManager()
	m.Create("g1")
	if err := m.Create("g1"); err == nil {
		t.Fatal("expected error for duplicate group")
	}
}

func TestCreateEmptyName(t *testing.T) {
	m := newManager()
	if err := m.Create(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAddAndMembers(t *testing.T) {
	m := newManager()
	m.Create("g1")
	m.Add("g1", "dev")
	m.Add("g1", "prod")
	m.Add("g1", "dev") // duplicate ignored
	members, err := m.Members("g1")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestRemoveMember(t *testing.T) {
	m := newManager()
	m.Create("g1")
	m.Add("g1", "dev")
	m.Add("g1", "prod")
	m.Remove("g1", "dev")
	members, _ := m.Members("g1")
	if len(members) != 1 || members[0] != "prod" {
		t.Fatalf("unexpected members: %v", members)
	}
}

func TestDeleteGroup(t *testing.T) {
	m := newManager()
	m.Create("g1")
	if err := m.Delete("g1"); err != nil {
		t.Fatal(err)
	}
	if err := m.Delete("g1"); err == nil {
		t.Fatal("expected error deleting non-existent group")
	}
}

func TestMembersNotFound(t *testing.T) {
	m := newManager()
	if _, err := m.Members("missing"); err == nil {
		t.Fatal("expected error")
	}
}
