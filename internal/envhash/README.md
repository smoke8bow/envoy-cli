# envhash

Package `envhash` provides deterministic SHA-256 hashing for environment variable maps.

## Overview

Hashing a profile lets you quickly detect whether its variables have changed
between two points in time — useful for cache invalidation, change detection,
and integrity checks.

## Usage

```go
h := envhash.New()

vars := map[string]string{
    "DATABASE_URL": "postgres://localhost/mydb",
    "PORT":         "8080",
}

hash := h.Compute(vars)
fmt.Println(hash) // deterministic hex string
```

## Key properties

- **Order-independent** — keys are sorted before hashing, so insertion order does not matter.
- **Deterministic** — the same map always produces the same hash.
- **Collision-resistant** — backed by SHA-256.

## API

| Function | Description |
|---|---|
| `New()` | Create a new `Hasher`. |
| `Compute(vars)` | Return the SHA-256 hex hash of a variable map. |
| `Equal(a, b)` | Report whether two maps hash identically. |
| `ComputeAll(profiles)` | Return sorted `[]Entry` for a map of named profiles. |
