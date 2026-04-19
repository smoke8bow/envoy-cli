package transform

import (
	"testing"
)

func TestIsSupportedValid(t *testing.T) {
	for _, op := range Supported() {
		if !IsSupported(op) {
			t.Errorf("expected %q to be supported", op)
		}
	}
}

func TestIsSupportedInvalid(t *testing.T) {
	if IsSupported(Op("rot13")) {
		t.Error("expected rot13 to be unsupported")
	}
}

func TestNewInvalidOp(t *testing.T) {
	_, err := New([]Op{"bad"})
	if err == nil {
		t.Fatal("expected error for unsupported op")
	}
}

func TestApplyUppercase(t *testing.T) {
	tr, _ := New([]Op{OpUppercase})
	out, err := tr.Apply(map[string]string{"KEY": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "HELLO" {
		t.Errorf("got %q", out["KEY"])
	}
}

func TestApplyLowercase(t *testing.T) {
	tr, _ := New([]Op{OpLowercase})
	out, _ := tr.Apply(map[string]string{"K": "WORLD"})
	if out["K"] != "world" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApplyTrimSpace(t *testing.T) {
	tr, _ := New([]Op{OpTrimSpace})
	out, _ := tr.Apply(map[string]string{"K": "  hi  "})
	if out["K"] != "hi" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApplyBase64RoundTrip(t *testing.T) {
	tr, _ := New([]Op{OpBase64Encode, OpBase64Decode})
	out, err := tr.Apply(map[string]string{"K": "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if out["K"] != "secret" {
		t.Errorf("got %q", out["K"])
	}
}

func TestApplyBase64DecodeInvalid(t *testing.T) {
	tr, _ := New([]Op{OpBase64Decode})
	_, err := tr.Apply(map[string]string{"K": "!!!not-base64!!!"})
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestApplyDoesNotMutateInput(t *testing.T) {
	tr, _ := New([]Op{OpUppercase})
	input := map[string]string{"K": "lower"}
	tr.Apply(input)
	if input["K"] != "lower" {
		t.Error("input was mutated")
	}
}
