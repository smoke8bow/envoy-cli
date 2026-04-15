# merge

The `merge` package provides utilities for combining two environment variable
maps with configurable conflict-resolution strategies.

## Strategies

| Strategy | Behaviour on conflict |
|---|---|
| `StrategyOurs` | Keep the destination value; record the key in `Skipped`. |
| `StrategyTheirs` | Overwrite with the source value; record the key in `Overwrite`. |
| `StrategyError` | Return an error immediately. |

## Usage

```go
result, err := merge.Merge(dst, src, merge.StrategyTheirs)
if err != nil {
    log.Fatal(err)
}
fmt.Println("added:", result.Added)
fmt.Println("overwritten:", result.Overwrite)
fmt.Println("final vars:", result.Vars)
```

The destination map is **never mutated**; a new map is always returned inside
`Result.Vars`.
