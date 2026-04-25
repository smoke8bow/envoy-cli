# envflag

`envflag` manages boolean feature flags stored inside a named environment
variable profile. Each flag is persisted as a `"true"` or `"false"` string
value, making it interoperable with any tool that reads the profile.

## Usage

```go
store  := // any Store implementation
m, err := envflag.NewManager(store, "myproject-flags")

// Enable a flag
_ = m.Set("DARK_MODE", true)

// Check a flag
enabled, err := m.IsEnabled("DARK_MODE")

// List all flags
flags, err := m.List() // map[string]bool

// Remove a flag
_ = m.Delete("DARK_MODE")
```

## Store accessor

When you already have a `ProfileStore` (the standard project store), use the
higher-level helper:

```go
m, err := envflag.NewStoreAccessor(profileStore, "myproject-flags")
```

This wraps the store so that a missing profile is treated as an empty flag
set rather than an error.

## Flag names

Flag names follow the same conventions as environment variable keys: they
should be non-empty. The package does not enforce upper-case naming but it is
recommended for consistency with the rest of the profile.
