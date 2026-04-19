package transform

import (
	"fmt"
	"strings"
)

// Op represents a transformation operation.
type Op string

const (
	OpUppercase Op = "uppercase"
	OpLowercase Op = "lowercase"
	OpTrimSpace Op = "trimspace"
	OpBase64Encode Op = "base64encode"
	OpBase64Decode Op = "base64decode"
)

// Supported returns all valid transformation ops.
func Supported() []Op {
	return []Op{OpUppercase, OpLowercase, OpTrimSpace, OpBase64Encode, OpBase64Decode}
}

// IsSupported returns true if op is a known transformation.
func IsSupported(op Op) bool {
	for _, s := range Supported() {
		if s == op {
			return true
		}
	}
	return false
}

// Transformer applies a sequence of Ops to env var values.
type Transformer struct {
	ops []Op
}

// New creates a Transformer with the given ops.
func New(ops []Op) (*Transformer, error) {
	for _, op := range ops {
		if !IsSupported(op) {
			return nil, fmt.Errorf("unsupported transform op: %q", op)
		}
	}
	return &Transformer{ops: ops}, nil
}

// Apply runs all ops on each value in vars and returns a new map.
func (t *Transformer) Apply(vars map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		result, err := t.applyValue(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = result
	}
	return out, nil
}

func (t *Transformer) applyValue(v string) (string, error) {
	for _, op := range t.ops {
		switch op {
		case OpUppercase:
			v = strings.ToUpper(v)
		case OpLowercase:
			v = strings.ToLower(v)
		case OpTrimSpace:
			v = strings.TrimSpace(v)
		case OpBase64Encode:
			v = b64Encode(v)
		case OpBase64Decode:
			dec, err := b64Decode(v)
			if err != nil {
				return "", fmt.Errorf("base64decode: %w", err)
			}
			v = dec
		}
	}
	return v, nil
}
