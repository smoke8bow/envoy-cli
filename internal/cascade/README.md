# cascade

The `cascade` package resolves environment variables across an **ordered chain of profiles**.

Profiles are applied left-to-right; later profiles override keys from earlier ones.
Each key in the final result is annotated with the name of the profile that last provided it.

## Usage

```go
acc := cascade.NewStoreAccessor(myStore.Get)
m   := cascade.NewManager(acc)

result, err := m.Resolve([]string{"base", "staging", "local"})
if err != nil {
    log.Fatal(err)
}

for _, k := range result.Keys() {
    fmt.Printf("%s=%s  (from %s)\n", k, result.Vars[k], result.Source[k])
}
```

## Behaviour

| Scenario | Result |
|---|---|
| Key only in first profile | kept as-is |
| Key in multiple profiles | last profile wins |
| Profile not found | error returned immediately |
| Empty profile list | error returned |
