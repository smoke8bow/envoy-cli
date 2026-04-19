# batch

The `batch` package provides atomic multi-operation updates to a single profile.

## Operations

| Kind     | Description                        |
|----------|------------------------------------|
| `set`    | Set or overwrite a key/value pair  |
| `delete` | Remove an existing key             |

## Usage

```go
proc := batch.NewProcessor(store)

ops := []batch.Op{
    {Key: "HOST", Value: "prod.example.com", Kind: batch.OpSet},
    {Key: "DEBUG", Kind: batch.OpDelete},
}

results, err := proc.Apply("production", ops)
```

## Behaviour

- All ops are validated before the profile is persisted.
- If any op fails the profile is left unchanged.
- Results contain per-op errors for inspection.
