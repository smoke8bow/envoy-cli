# freeze

The `freeze` package allows profiles to be marked as **frozen**, preventing any modifications to their environment variables.

## Usage

```go
m, err := freeze.NewManager("/path/to/frozen.json")
if err != nil {
    log.Fatal(err)
}

// Freeze a profile
_ = m.Freeze("prod")

// Check before mutating
if err := m.Guard("prod"); err != nil {
    fmt.Println("cannot modify:", err) // profile is frozen
}

// List all frozen profiles
for _, name := range m.List() {
    fmt.Println(name)
}

// Unfreeze when needed
_ = m.Unfreeze("prod")
```

## Errors

| Error | Meaning |
|---|---|
| `ErrFrozen` | The profile is frozen and cannot be modified |
| `ErrNotFrozen` | Unfreeze was called on a profile that is not frozen |

## Storage

Frozen state is persisted as a JSON file. The parent directory is created automatically if it does not exist.
