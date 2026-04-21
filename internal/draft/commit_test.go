package draft_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/user/envoy-cli/internal/draft"
)

type fakeWriter struct {
	saved map[string]map[string]string
	fail  bool
}

func (f *fakeWriter) Save(name string, vars map[string]string) error {
	if f.fail {
		return fmt.Errorf("write error")
	}
	if f.saved == nil {
		f.saved = make(map[string]map[string]string)
	}
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	f.saved[name] = copy
	return nil
}

func TestCommitUsesProfileName(t *testing.T) {
	m := draft.NewManager()
	_ = m.Create("tmp")
	_ = m.Set("tmp", "KEY", "val")
	w := &fakeWriter{}
	c := draft.NewCommitter(m, w)
	if err := c.Commit("tmp", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.saved["production"]["KEY"] != "val" {
		t.Errorf("expected KEY=val in production profile")
	}
	if _, err := m.Get("tmp"); !errors.Is(err, draft.ErrNoDraft) {
		t.Error("expected draft to be discarded after commit")
	}
}

func TestCommitFallsBackToDraftName(t *testing.T) {
	m := draft.NewManager()
	_ = m.Create("dev")
	_ = m.Set("dev", "A", "1")
	w := &fakeWriter{}
	c := draft.NewCommitter(m, w)
	_ = c.Commit("dev", "")
	if w.saved["dev"]["A"] != "1" {
		t.Error("expected draft name used as profile name")
	}
}

func TestCommitMissingDraft(t *testing.T) {
	m := draft.NewManager()
	w := &fakeWriter{}
	c := draft.NewCommitter(m, w)
	err := c.Commit("ghost", "prod")
	if !errors.Is(err, draft.ErrNoDraft) {
		t.Errorf("expected ErrNoDraft, got %v", err)
	}
}

func TestCommitWriterError(t *testing.T) {
	m := draft.NewManager()
	_ = m.Create("d")
	w := &fakeWriter{fail: true}
	c := draft.NewCommitter(m, w)
	if err := c.Commit("d", "p"); err == nil {
		t.Error("expected error from writer")
	}
}
