# envpromote

Promote environment variables from one named profile to another.

## Overview

`envpromote` copies keys from a **source** profile into a **destination** profile,
optionally filtering to a subset of keys and controlling whether existing
destination keys are overwritten.

## Usage

```go
store := // your Store implementation
m := envpromote.NewManager(store)

// Promote all keys from "dev" into "staging", overwriting conflicts.
result, err := m.Promote("dev", "staging", envpromote.DefaultOptions())

// Promote only specific keys without overwriting existing ones.
opts := envpromote.Options{
    Keys:      []string{"DB_HOST", "DB_PORT"},
    Overwrite: false,
}
result, err = m.Promote("dev", "prod", opts)
```

## Options

| Field | Default | Description |
|-------|---------|-------------|
| `Keys` | `[]` (all) | Restrict promotion to these keys |
| `Overwrite` | `true` | Replace existing keys in destination |

## Notes

- Source and destination must be different profiles.
- Missing keys in the source are silently skipped when `Keys` is specified.
- The destination profile is created if it does not exist (depends on the Store).
