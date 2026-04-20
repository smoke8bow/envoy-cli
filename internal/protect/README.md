# protect

The `protect` package allows marking specific environment variable keys within a profile as **protected**, preventing accidental modification or deletion.

## Usage

```go
m := protect.NewManager()

// Mark a key as protected
m.Protect("prod", "DB_PASSWORD")

// Check protection status
m.IsProtected("prod", "DB_PASSWORD") // true

// List all protected keys in a profile
keys := m.List("prod")

// Guard a write operation — returns ErrKeyProtected if any key is protected
err := m.Guard("prod", []string{"DB_PASSWORD", "HOST"})

// Remove protection
m.Unprotect("prod", "DB_PASSWORD")
```

## Accessor helpers

`GuardWrite` and `GuardDelete` wrap `Guard` with descriptive error context, suitable for use in store middleware.

## Errors

| Error | Meaning |
|---|---|
| `ErrKeyProtected` | Operation blocked because the key is protected |
| `ErrKeyNotProtected` | Unprotect called on a key that was never protected |
