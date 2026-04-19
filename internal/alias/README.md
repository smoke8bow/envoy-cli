# alias

The `alias` package lets you register short human-friendly names that map to
full profile names managed by envoy-cli.

## Usage

```go
store  := // any store.Store implementation
mgr    := alias.NewManager(store)

// register
_ = mgr.Set("prod", "production-us-east")

// resolve before switching
profile, _ := mgr.Resolve("prod") // "production-us-east"

// remove
_ = mgr.Remove("prod")

// list all
all := mgr.List() // map[string]string
```

## Rules

- Alias names must match `[a-zA-Z0-9_-]+`.
- The target profile must already exist in the store at the time `Set` is called.
- Aliases are held in memory; persist them via your own serialisation layer if needed.
