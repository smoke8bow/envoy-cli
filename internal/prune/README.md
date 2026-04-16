# prune

The `prune` package removes empty profiles from the store.

## Usage

```go
m := prune.NewManager(store)

// Preview what would be removed
candidates, err := m.DryRun()

// Actually remove empty profiles
result, err := m.Run()
fmt.Println("removed:", result.Removed)
fmt.Println("skipped:", result.Skipped)
```

## Behaviour

- A profile is considered **prunable** when its environment variable map is empty.
- `DryRun` lists candidates without mutating the store.
- `Run` deletes each candidate; profiles that fail to delete are added to `Skipped` rather than returning an error.
