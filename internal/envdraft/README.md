# envdraft

`envdraft` provides an in-memory staging area for environment variable edits
before they are committed to a named profile in the store.

## Concepts

- **Draft** – a temporary copy of a profile's variables that can be freely
  edited without affecting the live profile.
- **Manager** – holds all open drafts and coordinates Open / Set / Delete /
  Commit / Discard operations.

## Usage

```go
m := envdraft.NewManager(store)

// Open a draft seeded from the "dev" profile
d, err := m.Open("dev")

// Stage edits
_ = m.Set("dev", "LOG_LEVEL", "debug")
_ = m.Delete("dev", "LEGACY_FLAG")

// Inspect the staged state
d, _ = m.Get("dev")
fmt.Println(d.Vars)

// Persist to the store and close the draft
_ = m.Commit("dev")

// Or throw away all staged changes
_ = m.Discard("dev")
```

## Errors

| Error | Meaning |
|---|---|
| `ErrNoDraft` | No open draft for the requested profile |
| `ErrDraftExists` | A draft is already open for that profile |
