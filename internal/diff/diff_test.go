package diff

import (
	"strings"
	"testing"
)

func TestComputeNoChanges(t *testing.T) {
	from := map[string]string{"FOO": "bar", "BAZ": "qux"}
	to := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes := Compute(from, to)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestComputeAdd(t *testing.T) {
	from := map[string]string{}
	to := map[string]string{"NEW_VAR": "hello"}
	changes := Compute(from, to)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Op != OpAdd || changes[0].Key != "NEW_VAR" || changes[0].NewValue != "hello" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestComputeRemove(t *testing.T) {
	from := map[string]string{"OLD_VAR": "bye"}
	to := map[string]string{}
	changes := Compute(from, to)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Op != OpRemove || changes[0].Key != "OLD_VAR" || changes[0].OldValue != "bye" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestComputeUpdate(t *testing.T) {
	from := map[string]string{"HOST": "localhost"}
	to := map[string]string{"HOST": "production.example.com"}
	changes := Compute(from, to)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Op != OpUpdate || c.OldValue != "localhost" || c.NewValue != "production.example.com" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestComputeSorted(t *testing.T) {
	from := map[string]string{}
	to := map[string]string{"Z_VAR": "1", "A_VAR": "2", "M_VAR": "3"}
	changes := Compute(from, to)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	if changes[0].Key != "A_VAR" || changes[1].Key != "M_VAR" || changes[2].Key != "Z_VAR" {
		t.Errorf("changes not sorted: %v", changes)
	}
}

func TestFormatNoChanges(t *testing.T) {
	out := Format(nil)
	if out != "no changes" {
		t.Errorf("expected 'no changes', got %q", out)
	}
}

func TestFormatOutput(t *testing.T) {
	changes := []Change{
		{Key: "ADD_ME", NewValue: "yes", Op: OpAdd},
		{Key: "REMOVE_ME", OldValue: "old", Op: OpRemove},
		{Key: "UPDATE_ME", OldValue: "v1", NewValue: "v2", Op: OpUpdate},
	}
	out := Format(changes)
	if !strings.Contains(out, "+ ADD_ME=yes") {
		t.Errorf("missing add line in output: %s", out)
	}
	if !strings.Contains(out, "- REMOVE_ME=old") {
		t.Errorf("missing remove line in output: %s", out)
	}
	if !strings.Contains(out, "~ UPDATE_ME: v1 -> v2") {
		t.Errorf("missing update line in output: %s", out)
	}
}
