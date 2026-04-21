# patch

The `patch` package applies a sequence of typed operations to an environment
variable map, producing a new map without mutating the original.

## Operations

| Kind     | Fields required          | Description                          |
|----------|--------------------------|--------------------------------------|
| `set`    | `Key`, `Value`           | Create or overwrite a key            |
| `delete` | `Key`                    | Remove a key (no-op if absent)       |
| `rename` | `Key`, `NewKey`          | Rename a key, preserving its value   |

## Usage

```go
p := patch.New()

ops := []patch.Op{
    {Kind: patch.OpSet,    Key: "APP_ENV",  Value: "production"},
    {Kind: patch.OpDelete, Key: "DEBUG"},
    {Kind: patch.OpRename, Key: "DB_PASS",  NewKey: "DATABASE_PASSWORD"},
}

result, err := p.Apply(profile.Vars, ops)
```

Ops are applied in order, so later operations see the results of earlier ones.
The input map is never modified.
