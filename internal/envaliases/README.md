# envaliases

Provides short-alias mapping for environment variable keys within named profiles.

## Overview

`envaliases.Manager` lets you register human-friendly aliases for long
environment variable keys and expand them back to their canonical forms
before writing to a profile.

## Usage

```go
m := envaliases.NewManager()

// Register aliases for the "prod" profile
m.Set("prod", "db",    "DATABASE_URL")
m.Set("prod", "cache", "REDIS_URL")

// Resolve a single alias
key, err := m.Resolve("prod", "db") // → "DATABASE_URL"

// Expand a whole variable map
vars := map[string]string{
    "db":    "postgres://localhost/mydb",
    "OTHER": "unchanged",
}
expanded := m.Expand("prod", vars)
// expanded == {"DATABASE_URL": "postgres://localhost/mydb", "OTHER": "unchanged"}

// List all aliases (sorted)
aliases := m.List("prod") // → ["cache", "db"]

// Remove an alias
m.Remove("prod", "cache")
```

## Notes

- Aliases are scoped per profile; the same alias can exist in multiple profiles
  with different canonical keys.
- `Expand` does **not** mutate the input map.
