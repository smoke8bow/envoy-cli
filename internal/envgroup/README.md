# envgroup

Groups environment variables within a profile by their shared key prefix.

## Overview

Many projects use conventions like `DB_HOST`, `DB_PORT`, `APP_ENV` to namespace
variables. `envgroup` detects these prefixes automatically and returns a
structured `Result` containing named `Group` values and any leftover ungrouped
keys.

## Usage

```go
opts := envgroup.DefaultOptions()
// Require at least 2 keys to form a group (default).
// Use "_" as the separator (default).
// Optionally strip the prefix from keys inside the group.
opts.StripPrefix = true

result := envgroup.GroupBy(vars, opts)
for _, g := range result.Groups {
    fmt.Println(g.Prefix, g.Vars)
}
fmt.Println("ungrouped:", result.Ungrouped)
```

## Options

| Field         | Default | Description                                      |
|---------------|---------|--------------------------------------------------|
| `MinSize`     | `2`     | Minimum keys to form a group                     |
| `Separator`   | `_`     | Delimiter used to detect the prefix segment      |
| `StripPrefix` | `false` | Remove prefix+separator from keys inside a group |

## Store Accessor

```go
result, err := envgroup.GroupProfile(store, "production", opts)
```

Loads a profile from any `ProfileGetter` and groups its variables.
