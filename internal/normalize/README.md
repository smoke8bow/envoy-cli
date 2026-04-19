# normalize

The `normalize` package provides key normalization for environment variable maps.

## Strategies

| Strategy | Description |
|----------|-------------|
| `upper`  | Converts all keys to UPPER_CASE |
| `lower`  | Converts all keys to lower_case |
| `snake`  | Converts keys to UPPER_SNAKE_CASE (replaces `-` and spaces with `_`) |

## Usage

```go
n, err := normalize.New(normalize.StrategySnake)
if err != nil {
    log.Fatal(err)
}

normalized := n.Apply(map[string]string{
    "my-api-key": "abc123",
    "db host":    "localhost",
})
// Result: {"MY_API_KEY": "abc123", "DB_HOST": "localhost"}
```

## Notes

- `Apply` never mutates the input map.
- If two keys collide after normalization, the last one wins.
