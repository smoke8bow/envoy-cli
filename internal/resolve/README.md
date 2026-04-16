# resolve

The `resolve` package expands variable references within a profile's environment
variable map before the values are applied to a shell session.

## Syntax

Both `$VAR` and `${VAR}` forms are supported.

## Sources

References are resolved against:
1. Other keys within the **same profile**.
2. An optional **ambient** map (e.g. the current process environment).

## Cycle detection

If two or more keys reference each other, `Resolve` returns an error rather
than looping indefinitely.

## Usage

```go
r := resolve.NewResolver(ambientVars)
expanded, err := r.Resolve(profileVars)
```
