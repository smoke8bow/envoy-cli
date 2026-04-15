# template

The `template` package provides lightweight variable interpolation for environment variable values within `envoy-cli` profiles.

## Features

- Supports both `${VAR}` and `$VAR` reference styles.
- **`Expand`** – replaces known variables; unknown references are left unchanged.
- **`ExpandStrict`** – like `Expand` but returns an error listing every unresolved variable.
- **`ExpandMap`** – applies expansion to all values in a `map[string]string`.

## Resolvers

A `Resolver` is a simple `func(key string) (string, bool)` that the expander calls for each variable it encounters.

| Constructor | Description |
|---|---|
| `MapResolver(m)` | Backed by a static `map[string]string`. |
| `EnvResolver()` | Falls back to the host process environment (`os.LookupEnv`). |
| `ChainResolver(r...)` | Tries each resolver in order; returns the first match. |

## Example

```go
vars := map[string]string{
    "DB_HOST": "localhost",
    "DB_PORT": "5432",
    "DB_URL":  "postgres://${DB_HOST}:${DB_PORT}/mydb",
}

resolver := template.ChainResolver(
    template.MapResolver(vars),
    template.EnvResolver(),
)

expanded := template.ExpandMap(vars, resolver)
// expanded["DB_URL"] == "postgres://localhost:5432/mydb"
```
