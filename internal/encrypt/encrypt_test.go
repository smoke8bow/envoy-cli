package encrypt_test

import (
	"bytes"
	"testing"

	"envoy-cli/internal/encrypt"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	passphrase := "super-secret"
	original := []byte("KEY=value\nOTHER=123")

	ciphertext, err := encrypt.Encrypt(passphrase, original)
	if err != nil {
		t.Fatalf("Encrypt: unexpected error: %v", err)
	}

	if bytes.Equal(ciphertext, original) {
		t.Fatal("Encrypt: ciphertext must differ from plaintext")
	}

	plaintext, err := encrypt.Decrypt(passphrase, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: unexpected error: %v", err)
	}

	if !bytes.Equal(plaintext, original) {
		t.Fatalf("Decrypt: got %q, want %q", plaintext, original)
	}
}

func TestEncryptProducesUniqueNonces(t *testing.T) {
	passphrase := "nonce-test"
	plaintext := []byte("same plaintext")

	c1, err := encrypt.Encrypt(passphrase, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := encrypt.Encrypt(passphrase, plaintext)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(c1, c2) {
		t.Fatal("Encrypt: two calls with same input must not produce identical output")
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	ciphertext, err := encrypt.Encrypt("correct", []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = encrypt.Decrypt("wrong", ciphertext)
	if err == nil {
		t.Fatal("Decrypt: expected error for wrong passphrase, got nil")
	}
}

func TestDecryptTruncatedData(t *testing.T) {
	_, err := encrypt.Decrypt("passphrase", []byte{0x01, 0x02})
	if err == nil {
		t.Fatal("Decrypt: expected error for truncated data, got nil")
	}
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("Decrypt: expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecryptEmptyData(t *testing.T) {
	_, err := encrypt.Decrypt("passphrase", []byte{})
	if err != encrypt.ErrInvalidCiphertext {
		t.Fatalf("Decrypt: expected ErrInvalidCiphertext for empty data, got %v", err)
	}
}
