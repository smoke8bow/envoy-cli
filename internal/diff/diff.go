package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Change represents a single environment variable change.
type Change struct {
	Key      string
	OldValue string
	NewValue string
	Op       Op
}

// Op describes the type of change.
type Op int

const (
	OpAdd    Op = iota // variable added
	OpRemove           // variable removed
	OpUpdate           // variable value changed
)

func (o Op) String() string {
	switch o {
	case OpAdd:
		return "add"
	case OpRemove:
		return "remove"
	case OpUpdate:
		return "update"
	default:
		return "unknown"
	}
}

// Compute returns the list of changes needed to go from 'from' to 'to'.
func Compute(from, to map[string]string) []Change {
	var changes []Change

	// Check for additions and updates.
	for k, newVal := range to {
		if oldVal, exists := from[k]; !exists {
			changes = append(changes, Change{Key: k, NewValue: newVal, Op: OpAdd})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, OldValue: oldVal, NewValue: newVal, Op: OpUpdate})
		}
	}

	// Check for removals.
	for k, oldVal := range from {
		if _, exists := to[k]; !exists {
			changes = append(changes, Change{Key: k, OldValue: oldVal, Op: OpRemove})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return changes
}

// Format returns a human-readable summary of the changes.
func Format(changes []Change) string {
	if len(changes) == 0 {
		return "no changes"
	}
	var sb strings.Builder
	for _, c := range changes {
		switch c.Op {
		case OpAdd:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.NewValue)
		case OpRemove:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.OldValue)
		case OpUpdate:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
