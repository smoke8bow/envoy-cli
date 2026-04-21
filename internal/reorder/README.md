# reorder

The `reorder` package provides utilities for reordering keys within a named
environment profile.

## Usage

```go
r := reorder.NewReorderer(store)

// Persist a new key order
result, err := r.Apply("production", []string{"APP_ENV", "DB_URL", "PORT"})

// Preview without saving
keys, err := r.Preview("production", []string{"APP_ENV", "DB_URL"})
```

## Behaviour

- Keys listed in `order` that exist in the profile are placed first, in the
  given order.
- Keys that exist in the profile but are **not** listed are appended after the
  ordered keys.
- Keys listed in `order` that do **not** exist in the profile are ignored.
- `Preview` is identical to `Apply` but does not write back to the store.

## Error Handling

- If the named profile does not exist, both `Apply` and `Preview` return an
  `ErrProfileNotFound` error.
- If `order` is empty, the original key order from the profile is returned
  unchanged.
