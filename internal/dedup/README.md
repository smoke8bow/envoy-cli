# dedup

Detects and removes environment variable keys that share duplicate values within a profile.

## Strategies

| Strategy | Behaviour |
|---|---|
| `keep-first` | Among duplicate-value keys (sorted), keep the first and remove the rest |
| `keep-last` | Keep the last key, remove the others |
| `remove-all` | Remove every key that participates in a duplicate-value group |

## Usage

```go
d, err := dedup.New(dedup.StrategyKeepFirst)
if err != nil {
    log.Fatal(err)
}

cleaned := d.Apply(profile.Vars)
```

## Finding duplicates without removing

```go
results := dedup.Find(profile.Vars)
for _, r := range results {
    fmt.Printf("value %q shared by: %v\n", r.Value, r.Keys)
}
```
