# promote

The `promote` package copies environment variable keys from one named profile
into another — for example, promoting configuration from `staging` to `prod`.

## Usage

```go
p := promote.NewPromoter(store)

// Promote all keys, overwriting existing values
written, err := p.Promote("staging", "prod", promote.PromoteOptions{
    Overwrite: true,
})

// Promote only specific keys without overwriting
written, err = p.Promote("staging", "prod", promote.PromoteOptions{
    Keys:      []string{"API_KEY", "FEATURE_FLAG"},
    Overwrite: false,
})
```

## Behaviour

- If `Keys` is empty, all keys from the source profile are candidates.
- If `Overwrite` is `false`, keys that already exist in the destination are
  silently skipped.
- An error is returned if a requested key does not exist in the source profile.
- The destination profile is saved only when at least one key is written.
