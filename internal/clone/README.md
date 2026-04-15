# clone

The `clone` package provides profile duplication for **envoy-cli**.

## Usage

```go
m := profile.NewManager(s)
c := clone.NewCloner(m)

// Simple clone
err := c.Clone("production", "staging")

// Clone and override specific vars
err = c.CloneWithOverrides("production", "dev", map[string]string{
    "API_URL": "http://localhost:8080",
})
```

## Behaviour

| Scenario | Result |
|---|---|
| Source profile missing | error |
| Destination already exists | error |
| Destination name invalid | error |
| Successful clone | new profile with identical vars |
| Clone with overrides | new profile with merged vars |

Cloned vars are deep-copied; mutating the clone never affects the original profile.
