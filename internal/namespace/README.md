# namespace

The `namespace` package provides grouping of profiles under named scopes.

## Concepts

A **namespace** is a named collection of profile references. It does not copy
or own profiles — it simply tracks which profile names belong to a scope.

## Usage

```go
fs := namespace.NewFileStore("/path/to/data")
m := namespace.NewFSManager(fs)

// Create a namespace
m.Create("staging")

// Assign profiles
m.Assign("staging", "api-staging")
m.Assign("staging", "db-staging")

// List all namespaces
ns, _ := m.List()

// Remove a profile from a namespace
m.Unassign("staging", "db-staging")

// Delete a namespace
m.Delete("staging")
```

## Storage

Namespaces are persisted to `namespaces.json` in the configured data directory.
