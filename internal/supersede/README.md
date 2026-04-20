# supersede

The `supersede` package allows targeted key overrides from one profile onto another, without performing a full merge or replacement.

## Usage

```go
store := // any Store implementation
m := supersede.NewManager(store)

// Copy specific keys from "production" into "staging"
applied, err := m.Apply("staging", "production", []string{"DB_HOST", "DB_PORT"})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Applied keys:", applied)
```

## Behaviour

- If `keys` is **empty**, all keys from the source profile are copied into the destination.
- If `keys` is provided, only those keys are copied. Keys that do not exist in the source are silently skipped.
- The destination profile retains all keys that are not overridden.
- The source profile is never modified.

## Interface

```go
type Store interface {
    Get(name string) (map[string]string, error)
    Set(name string, vars map[string]string) error
}
```

Any store that satisfies this interface (e.g. `internal/store`) can be used.
