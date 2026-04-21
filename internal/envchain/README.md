# envchain

The `envchain` package provides **chained environment variable resolution** across multiple named profiles.

Profiles are layered in the order they are specified. Later profiles have higher priority and their values overwrite keys from earlier profiles.

## Usage

```go
chain, err := envchain.NewChain(profileStore, []string{"base", "production", "local-overrides"})
if err != nil {
    log.Fatal(err)
}

// Merge all profiles — later entries win on key conflicts.
resolved := chain.Resolve()
fmt.Println(resolved["DATABASE_URL"])

// Find which profile a key originates from.
origin := chain.Source("DATABASE_URL")
fmt.Println("comes from:", origin)
```

## Priority Order

| Position | Priority |
|----------|----------|
| First    | Lowest   |
| Last     | Highest  |

## Integration

`NewChain` accepts any value that satisfies the `ProfileGetter` interface:

```go
type ProfileGetter interface {
    Get(name string) (map[string]string, error)
}
```

The `profile.Manager` and `store.Store` types both satisfy this interface.
