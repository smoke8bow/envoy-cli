// Package patch provides fine-grained key-level updates to a profile's
// environment variables using a set of typed operations.
package patch

import "fmt"

// OpKind represents the type of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
	OpRename OpKind = "rename"
)

// Op describes a single patch operation.
type Op struct {
	Kind    OpKind
	Key     string
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Patcher applies a sequence of Ops to an env map.
type Patcher struct{}

// New returns a new Patcher.
func New() *Patcher { return &Patcher{} }

// Apply executes each Op in order against a copy of src and returns the result.
// The original map is never mutated.
func (p *Patcher) Apply(src map[string]string, ops []Op) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	for _, op := range ops {
		switch op.Kind {
		case OpSet:
			if op.Key == "" {
				return nil, fmt.Errorf("patch: set op missing key")
			}
			out[op.Key] = op.Value

		case OpDelete:
			if op.Key == "" {
				return nil, fmt.Errorf("patch: delete op missing key")
			}
			delete(out, op.Key)

		case OpRename:
			if op.Key == "" || op.NewKey == "" {
				return nil, fmt.Errorf("patch: rename op requires both key and new_key")
			}
			v, ok := out[op.Key]
			if !ok {
				return nil, fmt.Errorf("patch: rename source key %q not found", op.Key)
			}
			out[op.NewKey] = v
			delete(out, op.Key)

		default:
			return nil, fmt.Errorf("patch: unknown op kind %q", op.Kind)
		}
	}

	return out, nil
}
