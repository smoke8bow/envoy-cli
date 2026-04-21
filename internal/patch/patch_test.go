package patch

import (
	"testing"
)

func newPatcher() *Patcher { return New() }

func TestApplySet(t *testing.T) {
	p := newPatcher()
	src := map[string]string{"A": "1"}
	out, err := p.Apply(src, []Op{{Kind: OpSet, Key: "B", Value: "2"}})
	if err != nil {
		t.Fatal(err)
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
	if out["A"] != "1" {
		t.Error("original key A should be preserved")
	}
}

func TestApplyDelete(t *testing.T) {
	p := newPatcher()
	src := map[string]string{"A": "1", "B": "2"}
	out, err := p.Apply(src, []Op{{Kind: OpDelete, Key: "A"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["A"]; ok {
		t.Error("key A should have been deleted")
	}
	if out["B"] != "2" {
		t.Error("key B should remain")
	}
}

func TestApplyRename(t *testing.T) {
	p := newPatcher()
	src := map[string]string{"OLD": "val"}
	out, err := p.Apply(src, []Op{{Kind: OpRename, Key: "OLD", NewKey: "NEW"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["OLD"]; ok {
		t.Error("old key should be removed after rename")
	}
	if out["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %q", out["NEW"])
	}
}

func TestApplyRenameNotFound(t *testing.T) {
	p := newPatcher()
	_, err := p.Apply(map[string]string{}, []Op{{Kind: OpRename, Key: "MISSING", NewKey: "X"}})
	if err == nil {
		t.Error("expected error renaming missing key")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	p := newPatcher()
	src := map[string]string{"A": "1"}
	_, err := p.Apply(src, []Op{{Kind: OpSet, Key: "A", Value: "99"}})
	if err != nil {
		t.Fatal(err)
	}
	if src["A"] != "1" {
		t.Error("Apply must not mutate the source map")
	}
}

func TestApplyUnknownOpKind(t *testing.T) {
	p := newPatcher()
	_, err := p.Apply(map[string]string{}, []Op{{Kind: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown op kind")
	}
}

func TestApplySetMissingKey(t *testing.T) {
	p := newPatcher()
	_, err := p.Apply(map[string]string{}, []Op{{Kind: OpSet, Key: "", Value: "v"}})
	if err == nil {
		t.Error("expected error when set op has empty key")
	}
}
