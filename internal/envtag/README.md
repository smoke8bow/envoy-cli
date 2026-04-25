# envtag

Provides key-level tagging for environment variables within a named profile.

Tags are stored as JSON metadata alongside the profile store, keyed by a
reserved meta key so they don't pollute the env variable namespace.

## Usage

```go
store := // any Store implementation (e.g. wrapping internal/store)
m := envtag.NewManager(store)

// Attach tags to a key
m.Set("production", "DB_PASSWORD", []string{"secret", "database"})

// Retrieve tags for a key
tags, _ := m.Get("production", "DB_PASSWORD")
// => ["database", "secret"]  (always sorted)

// List all tagged keys in a profile
all, _ := m.List("production")
// => map[string][]string{"DB_PASSWORD": ["database", "secret"]}

// Remove tags from a key
m.Remove("production", "DB_PASSWORD")
```

## Store Interface

The `Store` interface requires only two methods:

```go
type Store interface {
    GetMeta(profile, key string) (string, error)
    SetMeta(profile, key, value string) error
}
```

This is satisfied by `internal/store` via its `meta.go` helpers.
