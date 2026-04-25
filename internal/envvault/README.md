# envvault

Provides encrypted storage for sensitive environment variable profiles.

Each profile's key-value map is serialised to JSON, encrypted with AES-256-GCM
(via the `encrypt` package), and stored as a base64 blob in a single JSON file.

## Usage

```go
fs := envvault.NewFileStorage("/home/user/.config/envoy/vault.json")
m  := envvault.NewManager(fs, passphrase)

// Store a new profile (fails if it already exists).
_ = m.Put("prod", map[string]string{"DB_PASS": "s3cr3t"})

// Overwrite an existing profile.
_ = m.Set("prod", map[string]string{"DB_PASS": "n3w-s3cr3t"})

// Retrieve and decrypt.
vars, _ := m.Get("prod")
fmt.Println(vars["DB_PASS"]) // n3w-s3cr3t

// List all vaulted profile names.
names, _ := m.List()

// Remove a profile.
_ = m.Delete("prod")
```

## Errors

| Sentinel | Meaning |
|---|---|
| `ErrNotFound` | Profile name does not exist in the vault |
| `ErrAlreadyExists` | `Put` called for a name that already exists |

## Security notes

- The vault file is written with mode `0600`.
- Each encryption call uses a random nonce, so identical plaintext produces
  different ciphertext on every write.
- The passphrase is never stored on disk.
