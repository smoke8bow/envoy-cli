package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config holds global CLI configuration settings.
type Config struct {
	DefaultShell    string `json:"default_shell,omitempty"`
	StorePath       string `json:"store_path,omitempty"`
	AuditEnabled    bool   `json:"audit_enabled"`
	EncryptionEnabled bool  `json:"encryption_enabled"`
	AutoExport      bool   `json:"auto_export"`
}

const configFileName = "config.json"

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		AuditEnabled:    true,
		EncryptionEnabled: false,
		AutoExport:      false,
	}
}

// Load reads the config file from dir, returning defaults if it does not exist.
func Load(dir string) (Config, error) {
	path := filepath.Join(dir, configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes the config to dir, creating the directory if necessary.
func Save(dir string, cfg Config) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, configFileName), data, 0o600)
}
