# envfence

The `envfence` package enforces an allowlist or denylist policy on the keys
present in an environment variable map.

## Modes

| Mode    | Behaviour                                              |
|---------|--------------------------------------------------------|
| `allow` | Only the listed keys are permitted; all others violate |
| `deny`  | The listed keys are forbidden; all others are allowed  |

## Usage

```go
// Allowlist – only HOST and PORT are permitted
f, err := envfence.New(envfence.ModeAllow, []string{"HOST", "PORT"})
if err != nil {
    log.Fatal(err)
}

vars := map[string]string{"HOST": "localhost", "SECRET": "abc"}

// Check returns violations without modifying vars
violations := f.Check(vars)
for _, v := range violations {
    fmt.Println(v) // key "SECRET": not in allowlist
}

// Filter returns a new map with disallowed keys removed
clean := f.Filter(vars) // {"HOST": "localhost"}
```

```go
// Denylist – SECRET is forbidden
f, _ := envfence.New(envfence.ModeDeny, []string{"SECRET"})
clean := f.Filter(vars)
```

## Notes

- `Check` and `Filter` never mutate the input map.
- Keys in the fence list are trimmed of surrounding whitespace on creation.
- At least one key must be supplied; an empty key list returns an error.
