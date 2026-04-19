# filter

The `filter` package provides key-based filtering of environment variable maps.

## Usage

```go
result := filter.Filter(vars, filter.Option{
    Prefix: "AWS_",
})
fmt.Println(result.Matched)  // only keys starting with AWS_
fmt.Println(result.Excluded) // all other keys
```

## Options

| Field | Description |
|-------|-------------|
| `Prefix` | Keep keys that start with this string |
| `Suffix` | Keep keys that end with this string |
| `Contains` | Keep keys that contain this substring |
| `ExactKeys` | Keep only these exact keys (overrides other options) |

Options are combined with AND logic (all non-empty options must match).
When `ExactKeys` is non-empty the other fields are ignored.
