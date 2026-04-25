package envsign_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/envsign"
)

func newSigner(t *testing.T, secret string) *envsign.Signer {
	t.Helper()
	s, err := envsign.New(secret)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNewEmptySecretError(t *testing.T) {
	_, err := envsign.New("")
	if err == nil {
		t.Fatal("expected error for empty secret")
	}
}

func TestSignDeterministic(t *testing.T) {
	s := newSigner(t, "mysecret")
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if s.Sign(vars) != s.Sign(vars) {
		t.Fatal("Sign should be deterministic")
	}
}

func TestSignOrderIndependent(t *testing.T) {
	s := newSigner(t, "mysecret")
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "2", "A": "1"}
	if s.Sign(a) != s.Sign(b) {
		t.Fatal("Sign should be order-independent")
	}
}

func TestVerifySuccess(t *testing.T) {
	s := newSigner(t, "secret")
	vars := map[string]string{"KEY": "value"}
	sig := s.Sign(vars)
	if err := s.Verify(vars, sig); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}

func TestVerifyTamperedValue(t *testing.T) {
	s := newSigner(t, "secret")
	vars := map[string]string{"KEY": "value"}
	sig := s.Sign(vars)
	vars["KEY"] = "tampered"
	if err := s.Verify(vars, sig); err == nil {
		t.Fatal("expected ErrInvalidSignature for tampered value")
	}
}

func TestVerifyTamperedKey(t *testing.T) {
	s := newSigner(t, "secret")
	vars := map[string]string{"ORIGINAL": "v"}
	sig := s.Sign(vars)
	delete(vars, "ORIGINAL")
	vars["MODIFIED"] = "v"
	if err := s.Verify(vars, sig); err == nil {
		t.Fatal("expected ErrInvalidSignature for tampered key")
	}
}

func TestVerifyWrongSecret(t *testing.T) {
	a := newSigner(t, "secret-a")
	b := newSigner(t, "secret-b")
	vars := map[string]string{"X": "y"}
	sig := a.Sign(vars)
	if err := b.Verify(vars, sig); err == nil {
		t.Fatal("expected ErrInvalidSignature for wrong secret")
	}
}

func TestSignEmptyMap(t *testing.T) {
	s := newSigner(t, "secret")
	sig := s.Sign(map[string]string{})
	if sig == "" {
		t.Fatal("expected non-empty signature for empty map")
	}
	if err := s.Verify(map[string]string{}, sig); err != nil {
		t.Fatalf("Verify empty map: %v", err)
	}
}
