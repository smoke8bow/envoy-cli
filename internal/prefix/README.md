# prefix

The `prefix` package provides utilities for applying, stripping, and filtering
environment variable key prefixes across a profile's variable map.

## Usage

```go
m := prefix.NewManager()

// Add a prefix to all keys
prefixed, err := m.Apply(vars, "MYAPP_")

// Remove the prefix from all keys (drop keys without the prefix)
stripped, err := m.Strip(prefixed, "MYAPP_", true)

// Keep only keys that carry the prefix
filtered := m.Filter(vars, "MYAPP_")
```

## Methods

| Method | Description |
|--------|-------------|
| `Apply(vars, prefix)` | Returns a new map with prefix prepended to every key |
| `Strip(vars, prefix, onlyPrefixed)` | Removes prefix from matching keys; optionally drops non-matching |
| `Filter(vars, prefix)` | Returns only entries whose keys start with prefix |
