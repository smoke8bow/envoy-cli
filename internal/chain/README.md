# chain

The `chain` package provides ordered composition of named environment variable layers.

## Overview

A `Composer` holds an ordered list of `Entry` values, each with a name and a
`map[string]string` of variables. When `Compose()` is called, entries are merged
left-to-right; later entries override keys from earlier ones.

The returned `Result` includes both the merged `Vars` map and a `Source` map that
tracks which layer each key was last set by.

## Usage

```go
entries := []chain.Entry{
    {Name: "base",     Vars: map[string]string{"DB_HOST": "localhost"}},
    {Name: "staging",  Vars: map[string]string{"DB_HOST": "staging.db", "LOG": "debug"}},
}
c, err := chain.NewComposer(entries)
if err != nil { ... }

result := c.Compose()
fmt.Println(result.Vars["DB_HOST"]) // staging.db
fmt.Println(result.Source["DB_HOST"]) // staging
```

## Loading from profiles

Use `FromProfiles` to build a `Composer` directly from a profile store:

```go
c, err := chain.FromProfiles(store, []string{"base", "staging"})
```
