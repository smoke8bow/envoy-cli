# envmap

The `envmap` package provides small, focused helpers for converting and
manipulating `map[string]string` environment variable collections.

## Functions

| Function | Description |
|---|---|
| `FromSlice(pairs []string)` | Parse a `KEY=VALUE` slice into a map |
| `ToSlice(m map[string]string)` | Serialize a map to a sorted `KEY=VALUE` slice |
| `Keys(m map[string]string)` | Return a sorted list of keys |
| `Clone(m map[string]string)` | Shallow-copy a map |
| `Merge(dst, src map[string]string)` | Merge src into dst (mutates dst) |

## Example

```go
m := envmap.FromSlice([]string{"FOO=bar", "BAZ=qux"})
slice := envmap.ToSlice(m)
// ["BAZ=qux", "FOO=bar"]

cloned := envmap.Clone(m)
envmap.Merge(cloned, map[string]string{"EXTRA": "1"})
```
