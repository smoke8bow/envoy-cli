package watch_test

import (
	"errors"
	"testing"

	"envoy-cli/internal/watch"
)

type fakeStore struct {
	data map[string]map[string]string
}

func (f *fakeStore) Get(name string) (map[string]string, error) {
	v, ok := f.data[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func newWatcher(data map[string]map[string]string) *watch.Watcher {
	return watch.NewWatcher(&fakeStore{data: data})
}

func TestChecksumDeterministic(t *testing.T) {
	vars := map[string]string{"B": "2", "A": "1"}
	if watch.Checksum(vars) != watch.Checksum(vars) {
		t.Fatal("expected same checksum")
	}
}

func TestChecksumDifferentValues(t *testing.T) {
	a := watch.Checksum(map[string]string{"KEY": "val1"})
	b := watch.Checksum(map[string]string{"KEY": "val2"})
	if a == b {
		t.Fatal("expected different checksums")
	}
}

func TestCheckBaseline(t *testing.T) {
	w := newWatcher(map[string]map[string]string{
		"dev": {"FOO": "bar"},
	})
	status, err := w.Check("dev", "")
	if err != nil {
		t.Fatal(err)
	}
	if status.Changed {
		t.Fatal("baseline should not be marked changed")
	}
	if status.Checksum == "" {
		t.Fatal("expected non-empty checksum")
	}
}

func TestCheckNoChange(t *testing.T) {
	vars := map[string]string{"FOO": "bar"}
	w := newWatcher(map[string]map[string]string{"dev": vars})
	baseline := watch.Checksum(vars)
	status, err := w.Check("dev", baseline)
	if err != nil {
		t.Fatal(err)
	}
	if status.Changed {
		t.Fatal("expected no change")
	}
}

func TestCheckDetectsChange(t *testing.T) {
	w := newWatcher(map[string]map[string]string{
		"dev": {"FOO": "new"},
	})
	old := watch.Checksum(map[string]string{"FOO": "old"})
	status, err := w.Check("dev", old)
	if err != nil {
		t.Fatal(err)
	}
	if !status.Changed {
		t.Fatal("expected change detected")
	}
}

func TestCheckProfileNotFound(t *testing.T) {
	w := newWatcher(map[string]map[string]string{})
	_, err := w.Check("missing", "")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
