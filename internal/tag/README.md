# tag

The `tag` package provides a lightweight tagging system for envoy-cli profiles.
Tags are stored as comma-separated metadata values inside the profile store,
keeping the implementation dependency-free and compatible with the existing
storage layer.

## Usage

```go
m := tag.NewManager(store)

// Attach tags
m.Add("dev", "backend")
m.Add("dev", "staging")

// List tags on a profile
tags, _ := m.List("dev") // ["backend", "staging"]

// Find all profiles with a specific tag
profiles, _ := m.ProfilesWithTag("backend")

// Remove a tag
m.Remove("dev", "staging")
```

## Storage

Tags are persisted via `store.SetMeta` / `store.GetMeta` under the key `"tags"`
as a comma-separated string. No additional files or databases are required.

## Constraints

- Tag names are trimmed of whitespace; empty tags are rejected.
- Duplicate tags on the same profile are silently ignored.
- Tags are always returned in alphabetical order.
