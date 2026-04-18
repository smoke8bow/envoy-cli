# archive

The `archive` package provides archiving and restoration of profile variable snapshots.

## Usage

```go
mgr := archive.NewManager("/path/to/archive/dir")

// Archive current state of a profile
err := mgr.Archive("production", map[string]string{"API_KEY": "abc"})

// List all archived entries for a profile
entries, err := mgr.List("production")

// Get the most recent archived entry
entry, err := mgr.Latest("production")
```

## Notes

- Each archive entry is stored as a timestamped JSON file.
- `Latest` returns the most recently archived entry for a given profile.
- `List` with an empty profile name returns all entries across all profiles.
