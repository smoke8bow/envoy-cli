package batch

import "fmt"

// Op represents a single batch operation on a profile's env vars.
type Op struct {
	Key   string
	Value string
	Kind  OpKind
}

// OpKind describes the type of operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
)

// Result holds the outcome of a single Op.
type Result struct {
	Op  Op
	Err error
}

// Store is the minimal interface required by the Processor.
type Store interface {
	Get(profile string) (map[string]string, error)
	Save(profile string, vars map[string]string) error
}

// Processor applies a list of Ops to a profile atomically.
type Processor struct {
	store Store
}

// NewProcessor creates a new Processor backed by the given Store.
func NewProcessor(s Store) *Processor {
	return &Processor{store: s}
}

// Apply executes all ops against the named profile and returns per-op results.
// The profile is only persisted if every op succeeds.
func (p *Processor) Apply(profile string, ops []Op) ([]Result, error) {
	vars, err := p.store.Get(profile)
	if err != nil {
		return nil, fmt.Errorf("batch: load profile %q: %w", profile, err)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	results := make([]Result, 0, len(ops))
	for _, op := range ops {
		var opErr error
		switch op.Kind {
		case OpSet:
			if op.Key == "" {
				opErr = fmt.Errorf("set: empty key")
			} else {
				copy[op.Key] = op.Value
			}
		case OpDelete:
			if _, ok := copy[op.Key]; !ok {
				opErr = fmt.Errorf("delete: key %q not found", op.Key)
			} else {
				delete(copy, op.Key)
			}
		default:
			opErr = fmt.Errorf("unknown op kind %q", op.Kind)
		}
		results = append(results, Result{Op: op, Err: opErr})
		if opErr != nil {
			return results, opErr
		}
	}

	if err := p.store.Save(profile, copy); err != nil {
		return results, fmt.Errorf("batch: save profile %q: %w", profile, err)
	}
	return results, nil
}
