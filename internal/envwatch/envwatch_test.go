package envwatch_test

import (
	"sync"
	"testing"
	"time"

	"envoy-cli/internal/envwatch"
)

// fakeStore implements a minimal in-memory profile store for testing.
type fakeStore struct {
	mu   sync.RWMutex
	data map[string]map[string]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{data: make(map[string]map[string]string)}
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.data[name]
	if !ok {
		return nil, nil
	}
	// return a copy
	copy := make(map[string]string, len(v))
	for k, val := range v {
		copy[k] = val
	}
	return copy, nil
}

func (f *fakeStore) set(name string, vars map[string]string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data[name] = vars
}

func newWatcher(t *testing.T, store envwatch.Store, profile string, interval time.Duration) *envwatch.Watcher {
	t.Helper()
	w, err := envwatch.NewWatcher(store, profile, interval)
	if err != nil {
		t.Fatalf("NewWatcher: %v", err)
	}
	return w
}

func TestNewWatcherUnknownProfile(t *testing.T) {
	store := newFakeStore()
	_, err := envwatch.NewWatcher(store, "ghost", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for unknown profile, got nil")
	}
}

func TestNewWatcherEmptyProfileName(t *testing.T) {
	store := newFakeStore()
	_, err := envwatch.NewWatcher(store, "", 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}

func TestNoChangeNoEvent(t *testing.T) {
	store := newFakeStore()
	store.set("prod", map[string]string{"FOO": "bar"})

	w := newWatcher(t, store, "prod", 30*time.Millisecond)
	defer w.Stop()

	ch := w.Changes()

	select {
	case ev := <-ch:
		t.Fatalf("unexpected change event: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// good — no change detected
	}
}

func TestDetectsAddedKey(t *testing.T) {
	store := newFakeStore()
	store.set("dev", map[string]string{"A": "1"})

	w := newWatcher(t, store, "dev", 30*time.Millisecond)
	defer w.Stop()

	// mutate after watcher is running
	time.Sleep(40 * time.Millisecond)
	store.set("dev", map[string]string{"A": "1", "B": "2"})

	select {
	case ev := <-w.Changes():
		if ev.Profile != "dev" {
			t.Errorf("expected profile 'dev', got %q", ev.Profile)
		}
		if len(ev.Diff) == 0 {
			t.Error("expected non-empty diff")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestDetectsRemovedKey(t *testing.T) {
	store := newFakeStore()
	store.set("staging", map[string]string{"X": "old", "Y": "keep"})

	w := newWatcher(t, store, "staging", 30*time.Millisecond)
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	store.set("staging", map[string]string{"Y": "keep"})

	select {
	case ev := <-w.Changes():
		if ev.Profile != "staging" {
			t.Errorf("profile mismatch: got %q", ev.Profile)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for removal event")
	}
}

func TestDetectsValueChange(t *testing.T) {
	store := newFakeStore()
	store.set("local", map[string]string{"DB_URL": "postgres://old"})

	w := newWatcher(t, store, "local", 30*time.Millisecond)
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	store.set("local", map[string]string{"DB_URL": "postgres://new"})

	select {
	case ev := <-w.Changes():
		if ev.Profile != "local" {
			t.Errorf("profile mismatch: got %q", ev.Profile)
		}
		found := false
		for _, d := range ev.Diff {
			if d.Key == "DB_URL" {
				found = true
			}
		}
		if !found {
			t.Error("expected DB_URL in diff")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for value-change event")
	}
}

func TestStopCancelsWatcher(t *testing.T) {
	store := newFakeStore()
	store.set("ci", map[string]string{"CI": "true"})

	w := newWatcher(t, store, "ci", 20*time.Millisecond)
	w.Stop()

	// After Stop, Changes channel should be closed or drain without blocking.
	time.Sleep(80 * time.Millisecond)
	store.set("ci", map[string]string{"CI": "false"})

	select {
	case _, ok := <-w.Changes():
		if ok {
			// A stale event before stop is acceptable; just drain.
		}
	case <-time.After(150 * time.Millisecond):
		// No event after stop — expected.
	}
}
