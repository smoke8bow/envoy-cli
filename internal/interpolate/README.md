# interpolate

The `interpolate` package resolves variable references embedded in profile
values. It supports both `${VAR}` and `$VAR` syntax and performs multiple
passes until all references are resolved or the depth limit is reached.

## Usage

```go
interp := interpolate.New(profileGetter, interpolate.DefaultOptions())
resolved, err := interp.Apply(vars, os.Getenv)
```

## Cross-profile references

References are resolved against the profile's own variables first. When
`FallbackToOS` is enabled (default), unresolved keys are looked up in the
OS environment.

## Options

| Field          | Default | Description                                      |
|----------------|---------|--------------------------------------------------|
| FallbackToOS   | true    | Fall back to `os.Getenv` for unresolved keys     |
| MaxDepth       | 8       | Maximum interpolation passes to prevent cycles   |
