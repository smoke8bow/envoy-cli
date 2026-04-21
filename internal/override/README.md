# override

The `override` package provides ephemeral, in-memory key-value overlays for
named profiles. Overrides are layered on top of a profile's stored variables
at runtime and are never written back to the store.

## Usage

```go
m := override.NewManager()

// Add an ephemeral override
m.Set("dev", "DATABASE_URL", "postgres://localhost/dev_tmp")

// Merge overrides on top of a profile's base vars
result := m.Apply("dev", profileVars)

// Inspect the current layer
layer := m.Layer("dev")

// Remove a single key from the layer
m.Unset("dev", "DATABASE_URL")

// Wipe the entire layer for a profile
m.Clear("dev")
```

## Behaviour

- `Apply` never mutates the `base` map; it always returns a new map.
- `Layer` returns a copy of the internal map so callers cannot accidentally
  mutate the manager's state.
- Overrides exist only for the lifetime of the process; they are not persisted
  to disk.
