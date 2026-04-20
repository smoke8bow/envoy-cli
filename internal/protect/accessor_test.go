package protect_test

import (
	"errors"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/protect"
)

type fakeStore struct {
	vars map[string]map[string]string
}

func (f *fakeStore) Get(profile string) (map[string]string, error) {
	v, ok := f.vars[profile]
	if !ok {
		return nil, errors.New("profile not found")
	}
	return v, nil
}

func TestGuardWriteAllowed(t *testing.T) {
	m := protect.NewManager()
	_ = m.Protect("prod", "LOCKED")
	store := &fakeStore{}
	err := protect.GuardWrite(m, store, "prod", []string{"SAFE_KEY"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGuardWriteBlocked(t *testing.T) {
	m := protect.NewManager()
	_ = m.Protect("prod", "LOCKED")
	store := &fakeStore{}
	err := protect.GuardWrite(m, store, "prod", []string{"LOCKED"})
	if !errors.Is(err, protect.ErrKeyProtected) {
		t.Errorf("expected ErrKeyProtected, got %v", err)
	}
}

func TestGuardDeleteAllowed(t *testing.T) {
	m := protect.NewManager()
	_ = m.Protect("prod", "LOCKED")
	err := protect.GuardDelete(m, "prod", []string{"OTHER"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGuardDeleteBlocked(t *testing.T) {
	m := protect.NewManager()
	_ = m.Protect("prod", "LOCKED")
	err := protect.GuardDelete(m, "prod", []string{"LOCKED", "SAFE"})
	if !errors.Is(err, protect.ErrKeyProtected) {
		t.Errorf("expected ErrKeyProtected, got %v", err)
	}
}
