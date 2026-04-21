package envdraft_test

import (
	"errors"
	"testing"

	"github.com/your-org/envoy-cli/internal/envdraft"
)

// fakeStore is a simple in-memory Store for testing.
type fakeStore struct {
	data map[string]map[string]string
}

func newFakeStore(profiles ...string) *fakeStore {
	s := &fakeStore{data: make(map[string]map[string]string)}
	for _, p := range profiles {
		s.data[p] = make(map[string]string)
	}
	return s
}

func (f *fakeStore) Get(profile string) (map[string]string, error) {
	v, ok := f.data[profile]
	if !ok {
		return nil, errors.New("profile not found")
	}
	return v, nil
}

func (f *fakeStore) Save(profile string, vars map[string]string) error {
	f.data[profile] = vars
	return nil
}

func newManager(t *testing.T) (*envdraft.Manager, *fakeStore) {
	t.Helper()
	s := newFakeStore("dev", "prod")
	s.data["dev"]["APP_ENV"] = "development"
	return envdraft.NewManager(s), s
}

func TestOpenCreatesDraft(t *testing.T) {
	m, _ := newManager(t)
	d, err := m.Open("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Vars["APP_ENV"] != "development" {
		t.Errorf("expected seeded value, got %q", d.Vars["APP_ENV"])
	}
}

func TestOpenDuplicateReturnsError(t *testing.T) {
	m, _ := newManager(t)
	_, _ = m.Open("dev")
	_, err := m.Open("dev")
	if !errors.Is(err, envdraft.ErrDraftExists) {
		t.Fatalf("expected ErrDraftExists, got %v", err)
	}
}

func TestSetAndGet(t *testing.T) {
	m, _ := newManager(t)
	_, _ = m.Open("dev")
	if err := m.Set("dev", "NEW_KEY", "hello"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	d, _ := m.Get("dev")
	if d.Vars["NEW_KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", d.Vars["NEW_KEY"])
	}
}

func TestDeleteKey(t *testing.T) {
	m, _ := newManager(t)
	_, _ = m.Open("dev")
	_ = m.Delete("dev", "APP_ENV")
	d, _ := m.Get("dev")
	if _, ok := d.Vars["APP_ENV"]; ok {
		t.Error("expected key to be deleted")
	}
}

func TestCommitPersists(t *testing.T) {
	m, s := newManager(t)
	_, _ = m.Open("dev")
	_ = m.Set("dev", "COMMITTED", "yes")
	if err := m.Commit("dev"); err != nil {
		t.Fatalf("Commit: %v", err)
	}
	if s.data["dev"]["COMMITTED"] != "yes" {
		t.Error("expected committed value in store")
	}
	// draft should be closed after commit
	if _, err := m.Get("dev"); !errors.Is(err, envdraft.ErrNoDraft) {
		t.Error("expected ErrNoDraft after commit")
	}
}

func TestDiscard(t *testing.T) {
	m, _ := newManager(t)
	_, _ = m.Open("dev")
	if err := m.Discard("dev"); err != nil {
		t.Fatalf("Discard: %v", err)
	}
	if _, err := m.Get("dev"); !errors.Is(err, envdraft.ErrNoDraft) {
		t.Error("expected ErrNoDraft after discard")
	}
}

func TestSetOnClosedDraft(t *testing.T) {
	m, _ := newManager(t)
	if err := m.Set("dev", "K", "v"); !errors.Is(err, envdraft.ErrNoDraft) {
		t.Errorf("expected ErrNoDraft, got %v", err)
	}
}

func TestOpenUnknownProfileError(t *testing.T) {
	m, _ := newManager(t)
	_, err := m.Open("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}
