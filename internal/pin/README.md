# pin

The `pin` package allows profiles to be marked as **pinned**, protecting them from accidental modification or deletion.

## Usage

```go
m := pin.NewManager(store)

// Pin a profile
m.Pin("production")

// Check if pinned
if m.IsPinned("production") {
    fmt.Println("profile is pinned")
}

// List all pinned profiles
names, _ := m.ListPinned()

// Unpin
m.Unpin("production")
```

## Notes

- Pinned state is stored in profile metadata under the key `pinned`.
- Callers responsible for enforcing pin guards before mutating profiles.
