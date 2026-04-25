# envrotate

Key rotation for environment variable profiles.

## Overview

`envrotate` renames one or more keys within a named profile, preserving their
values. This is useful when standardising key naming conventions across
projects (e.g. renaming `DB_HOST` → `DATABASE_HOST`).

## Usage

```go
store := envrotate.NewStoreAccessor(myStore.Get, myStore.Set)
m := envrotate.NewManager(store, envrotate.DefaultOptions())

result, err := m.Rotate("production", map[string]string{
    "DB_HOST": "DATABASE_HOST",
    "DB_PASS": "DATABASE_PASSWORD",
})
if err != nil {
    log.Fatal(err)
}

for _, r := range result.Rotated {
    fmt.Printf("renamed %s -> %s\n", r.OldKey, r.NewKey)
}
for _, s := range result.Skipped {
    fmt.Printf("skipped (not found): %s\n", s)
}
```

## Options

| Field       | Default | Description                                      |
|-------------|---------|--------------------------------------------------|
| `RemoveOld` | `true`  | Delete the old key after copying its value over. |

Set `RemoveOld: false` to copy the value to the new key while keeping the
original key intact — useful for a phased migration.

## Convenience

`RotateProfile` performs a one-shot rotation without constructing a Manager:

```go
result, err := envrotate.RotateProfile(store, "staging", rotMap, opts)
```
