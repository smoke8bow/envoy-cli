# group

The `group` package provides named collections of profiles.

## Usage

```go
store := myStore // implements group.Store
m := group.NewManager(store)

// Create a group
m.Create("staging")

// Add profiles to the group
m.Add("staging", "staging-us")
m.Add("staging", "staging-eu")

// List members
members, _ := m.Members("staging")
// => ["staging-eu", "staging-us"]

// Remove a profile
m.Remove("staging", "staging-eu")

// List all groups
groups := m.List()

// Delete a group
m.Delete("staging")
```

## Concepts

- Groups are in-memory collections; persistence is handled by the caller.
- Adding a duplicate profile to a group is a no-op.
- All list operations return sorted results for deterministic output.
