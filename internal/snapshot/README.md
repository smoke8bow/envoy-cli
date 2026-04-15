# snapshot

The `snapshot` package provides point-in-time captures of a profile's environment variables.

## Overview

A **snapshot** is an immutable copy of a named profile saved under a generated ID:

```
<profile>__snap_<unix-milliseconds>
```

Snapshots are stored as regular profiles in the backing store, so they survive restarts and can be listed, exported, or applied like any other profile.

## Usage

```go
m := snapshot.NewManagerFromStore(s)

// Capture current state of "prod"
entry, err := m.Take("prod", "before v2 deploy")

// Restore an earlier snapshot back to "prod"
err = m.Restore(entry.ID, "prod")
```

## Entry fields

| Field       | Description                              |
|-------------|------------------------------------------|
| `ID`        | Unique snapshot identifier               |
| `Profile`   | Source profile name                      |
| `Vars`      | Copy of environment variables at capture |
| `CreatedAt` | UTC timestamp of the snapshot            |
| `Note`      | Optional human-readable label            |
