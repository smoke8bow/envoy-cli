# envdiff

The `envdiff` package compares a named profile's environment variables against the current OS environment, showing what is new, shared, or only in the OS.

## Usage

```go
result := envdiff.Compare(profileVars, os.Environ())
fmt.Println(envdiff.Format(result))
```

## Result Fields

Each entry in the result slice carries:

| Field | Description |
|-------|-------------|
| `Key` | Environment variable name |
| `ProfileValue` | Value from the named profile (empty if not present) |
| `OSValue` | Value from the running OS environment (empty if not present) |
| `Source` | One of `profile-only`, `os-only`, or `both` |

## Format Output

`Format` renders a human-readable diff where:

- `+` prefix — key exists only in the profile
- `-` prefix — key exists only in the OS environment
- `~` prefix — key exists in both but values differ
- ` ` (space) prefix — key exists in both with identical values

## Notes

- Results are sorted alphabetically by key for deterministic output.
- `Compare` does not modify the input slices.
