# defaults

The `defaults` package lets you register fallback environment variable values
for a named profile. When a profile is applied, any key that is **not** already
present in the active variable set is filled in from the registered defaults.

## Usage

```go
m := defaults.NewManager("/path/to/data/dir")

// Register defaults for a profile
m.Set("dev", map[string]string{
    "LOG_LEVEL": "debug",
    "TIMEOUT":   "30",
})

// Retrieve defaults
vars, _ := m.Get("dev")

// Merge defaults into an existing variable map (existing keys win)
result, _ := m.Apply("dev", currentVars)

// Remove defaults for a profile
m.Delete("dev")
```

## Behaviour

- `Apply` never overwrites keys that already exist in the supplied map.
- `Get` returns an empty map (no error) when no defaults have been registered.
- State is persisted to `defaults.json` inside the directory passed to
  `NewManager`, so it survives process restarts.
