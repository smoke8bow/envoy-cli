# watch

The `watch` package provides checksum-based change detection for environment profiles.

## Usage

```go
w := watch.NewWatcher(profileStore)

// Baseline — record current checksum
status, err := w.Check("production", "")
baseline := status.Checksum

// Later — detect if profile changed
status, err = w.Check("production", baseline)
if status.Changed {
    fmt.Println("Profile has changed since last check")
}
```

## Checksum

`Checksum(vars)` produces a deterministic SHA-256 hex digest of an env map by
sorting keys before hashing, ensuring consistent results regardless of map
iteration order.
