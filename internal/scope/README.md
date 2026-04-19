# scope

The `scope` package provides a lightweight named-context system for grouping and filtering profiles by arbitrary label selectors.

## Usage

```go
m := scope.NewManager()

// Create a scope with labels
s, err := m.Create("prod", map[string]string{"env": "production", "region": "us-east"})

// Retrieve a scope
s, err = m.Get("prod")

// Match against selectors
if s.Match(map[string]string{"env": "production"}) {
    // apply prod-specific behaviour
}

// List all scopes
scopes := m.List()

// Delete a scope
err = m.Delete("prod")
```

## Labels

Labels are arbitrary key-value pairs attached to a scope. Use `Match` to check whether a scope satisfies a set of selector requirements (all selectors must match).
