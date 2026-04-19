# required

The `required` package enforces that specific environment variable keys are present and non-empty within a named profile.

## Usage

```go
m := required.NewManager()
m.Set("production", []string{"DATABASE_URL", "SECRET_KEY"})

violations := m.Check("production", map[string]string{
    "DATABASE_URL": "postgres://...",
})
// violations: [{Key: "SECRET_KEY", Profile: "production"}]
```

## API

- `NewManager()` — create a new manager
- `Set(profile, keys)` — define required keys for a profile
- `Get(profile)` — retrieve required keys for a profile
- `Check(profile, vars)` — validate vars against requirements, returns violations
- `CheckAll(fn)` — validate all registered profiles using a vars-fetching function
