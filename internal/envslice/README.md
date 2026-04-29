# envslice

Provides utilities for converting environment variable maps to ordered slices
of `Entry` values and back.

## Types

```go
type Entry struct {
    Key   string
    Value string
}
```

## Functions

| Function | Description |
|---|---|
| `FromMap(vars map[string]string) []Entry` | Convert map to sorted entry slice |
| `ToMap(entries []Entry) map[string]string` | Convert entries back to map |
| `ToStrings(entries []Entry) []string` | Render entries as `KEY=VALUE` strings |
| `FromStrings(lines []string) []Entry` | Parse `KEY=VALUE` strings into entries |
| `FilterByPrefix(entries []Entry, prefix string) []Entry` | Keep only entries whose key starts with prefix |

## Example

```go
vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
entries := envslice.FromMap(vars)
appEntries := envslice.FilterByPrefix(entries, "APP_")
for _, e := range appEntries {
    fmt.Println(e) // APP_HOST=localhost, APP_PORT=8080
}
```
