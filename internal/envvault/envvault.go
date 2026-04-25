// Package envvault provides encrypted storage for sensitive environment profiles.
package envvault

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/user/envoy-cli/internal/encrypt"
)

// ErrNotFound is returned when a vaulted profile does not exist.
var ErrNotFound = errors.New("vaulted profile not found")

// ErrAlreadyExists is returned when a vaulted profile already exists.
var ErrAlreadyExists = errors.New("vaulted profile already exists")

// Storage is the interface required for persisting encrypted blobs.
type Storage interface {
	Load() (map[string]string, error)
	Save(map[string]string) error
}

// Manager manages encrypted environment profiles in a vault.
type Manager struct {
	storage    Storage
	passphrase string
}

// NewManager returns a Manager backed by the given Storage and passphrase.
func NewManager(storage Storage, passphrase string) *Manager {
	return &Manager{storage: storage, passphrase: passphrase}
}

// Put encrypts and stores vars under name, returning ErrAlreadyExists if it exists.
func (m *Manager) Put(name string, vars map[string]string) error {
	blobs, err := m.storage.Load()
	if err != nil {
		return fmt.Errorf("envvault: load: %w", err)
	}
	if _, ok := blobs[name]; ok {
		return fmt.Errorf("%w: %s", ErrAlreadyExists, name)
	}
	return m.write(blobs, name, vars)
}

// Set encrypts and stores vars under name, overwriting any existing entry.
func (m *Manager) Set(name string, vars map[string]string) error {
	blobs, err := m.storage.Load()
	if err != nil {
		return fmt.Errorf("envvault: load: %w", err)
	}
	return m.write(blobs, name, vars)
}

// Get decrypts and returns the vars stored under name.
func (m *Manager) Get(name string) (map[string]string, error) {
	blobs, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("envvault: load: %w", err)
	}
	blob, ok := blobs[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	plain, err := encrypt.Decrypt(blob, m.passphrase)
	if err != nil {
		return nil, fmt.Errorf("envvault: decrypt %s: %w", name, err)
	}
	var vars map[string]string
	if err := json.Unmarshal([]byte(plain), &vars); err != nil {
		return nil, fmt.Errorf("envvault: unmarshal %s: %w", name, err)
	}
	return vars, nil
}

// Delete removes a vaulted profile by name.
func (m *Manager) Delete(name string) error {
	blobs, err := m.storage.Load()
	if err != nil {
		return fmt.Errorf("envvault: load: %w", err)
	}
	if _, ok := blobs[name]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	delete(blobs, name)
	return m.storage.Save(blobs)
}

// List returns all vaulted profile names.
func (m *Manager) List() ([]string, error) {
	blobs, err := m.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("envvault: load: %w", err)
	}
	names := make([]string, 0, len(blobs))
	for k := range blobs {
		names = append(names, k)
	}
	return names, nil
}

func (m *Manager) write(blobs map[string]string, name string, vars map[string]string) error {
	raw, err := json.Marshal(vars)
	if err != nil {
		return fmt.Errorf("envvault: marshal: %w", err)
	}
	cipher, err := encrypt.Encrypt(string(raw), m.passphrase)
	if err != nil {
		return fmt.Errorf("envvault: encrypt: %w", err)
	}
	blobs[name] = cipher
	return m.storage.Save(blobs)
}
