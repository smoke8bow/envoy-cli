package envtag_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/envtag"
)

// fakeStore is an in-memory Store for testing.
type fakeStore struct {
	meta map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{meta: map[string]string{}}
}

func (f *fakeStore) GetMeta(profile, key string) (string, error) {
	return f.meta[profile+":"+key], nil
}

func (f *fakeStore) SetMeta(profile, key, value string) error {
	f.meta[profile+":"+key] = value
	return nil
}

func newManager(t *testing.T) *envtag.Manager {
	t.Helper()
	return envtag.NewManager(newFakeStore())
}

func TestSetAndGet(t *testing.T) {
	m := newManager(t)
	if err := m.Set("prod", "DB_URL", []string{"secret", "db"}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	tags, err := m.Get("prod", "DB_URL")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestGetUntaggedKey(t *testing.T) {
	m := newManager(t)
	tags, err := m.Get("prod", "MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %v", tags)
	}
}

func TestRemoveTag(t *testing.T) {
	m := newManager(t)
	_ = m.Set("prod", "API_KEY", []string{"secret"})
	if err := m.Remove("prod", "API_KEY"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	tags, _ := m.Get("prod", "API_KEY")
	if len(tags) != 0 {
		t.Fatalf("expected empty after remove, got %v", tags)
	}
}

func TestListReturnsAll(t *testing.T) {
	m := newManager(t)
	_ = m.Set("dev", "FOO", []string{"a"})
	_ = m.Set("dev", "BAR", []string{"b", "c"})
	all, err := m.List("dev")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestTagsSorted(t *testing.T) {
	m := newManager(t)
	_ = m.Set("prod", "KEY", []string{"z", "a", "m"})
	tags, _ := m.Get("prod", "KEY")
	if tags[0] != "a" || tags[1] != "m" || tags[2] != "z" {
		t.Fatalf("expected sorted tags, got %v", tags)
	}
}

func TestSetEmptyProfileError(t *testing.T) {
	m := newManager(t)
	if err := m.Set("", "KEY", []string{"tag"}); err == nil {
		t.Fatal("expected error for empty profile")
	}
}

func TestSetEmptyKeyError(t *testing.T) {
	m := newManager(t)
	if err := m.Set("prod", "", []string{"tag"}); err == nil {
		t.Fatal("expected error for empty key")
	}
}
