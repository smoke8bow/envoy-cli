# pivot

The `pivot` package provides a utility to transpose a profile's environment
variable map — swapping keys and values.

## Directions

| Direction | Behaviour |
|---|---|
| `keys-to-values` | Each key becomes a value and each value becomes a key |
| `values-to-keys` | Alias for the same swap operation |

## Usage

```go
result, err := pivot.Pivot(vars, pivot.DirectionKeysToValues)
```

An error is returned when duplicate values exist, as they cannot be safely
promoted to unique keys.
