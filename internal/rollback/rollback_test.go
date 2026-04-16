package rollback

import (
	"errors"
	"testing"
)

// --- fakes ---

type fakeStore struct {
	profiles map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.profiles[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (f *fakeStore) Set(name string, vars map[string]string) error {
	f.profiles[name] = vars
	return nil
}

type fakeSnapshot struct {
	snaps    map[string]string // name -> profile
	takeErr  error
	restored string
	restoreErr error
}

func (f *fakeSnapshot) Take(profile string) (string, error) {
	if f.takeErr != nil {
		return "", f.takeErr
	}
	name := "snap-" + profile
	f.snaps[name] = profile
	return name, nil
}

func (f *fakeSnapshot) Restore(snapshotName string) error {
	if f.restoreErr != nil {
		return f.restoreErr
	}
	f.restored = snapshotName
	return nil
}

func newManager() (*Manager, *fakeStore, *fakeSnapshot) {
	st := &fakeStore{profiles: map[string]map[string]string{
		"dev": {"KEY": "val"},
	}}
	sn := &fakeSnapshot{snaps: map[string]string{}}
	return NewManager(st, sn), st, sn
}

func TestCheckpointSuccess(t *testing.T) {
	m, _, sn := newManager()
	name, err := m.Checkpoint("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "snap-dev" {
		t.Errorf("expected snap-dev, got %s", name)
	}
	if sn.snaps[name] != "dev" {
		t.Errorf("snapshot not recorded")
	}
}

func TestCheckpointProfileNotFound(t *testing.T) {
	m, _, _ := newManager()
	_, err := m.Checkpoint("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestRollbackSuccess(t *testing.T) {
	m, _, sn := newManager()
	if err := m.Rollback("snap-dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sn.restored != "snap-dev" {
		t.Errorf("expected snap-dev to be restored")
	}
}

func TestRollbackEmptyName(t *testing.T) {
	m, _, _ := newManager()
	if err := m.Rollback(""); err == nil {
		t.Fatal("expected error for empty snapshot name")
	}
}

func TestRollbackRestoreError(t *testing.T) {
	m, _, sn := newManager()
	sn.restoreErr = errors.New("disk full")
	if err := m.Rollback("snap-dev"); err == nil {
		t.Fatal("expected error from restore")
	}
}
