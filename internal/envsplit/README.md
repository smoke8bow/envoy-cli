# envsplit

Splits a flat environment variable map into named **buckets** based on key-prefix rules.

## Usage

```go
rules := []envsplit.Rule{
    {Prefix: "APP_", Bucket: "app"},
    {Prefix: "DB_",  Bucket: "db"},
}

opts := envsplit.Options{StripPrefix: true}

result, err := envsplit.Split(src, rules, opts)
// result.Buckets["app"] → keys with APP_ prefix (prefix stripped)
// result.Buckets["db"]  → keys with DB_ prefix  (prefix stripped)
// result.Remainder      → unmatched keys
```

## Rules

- Rules are evaluated **in order**; the **first match wins**.
- `Prefix` and `Bucket` must both be non-empty.
- When `StripPrefix` is `true`, the matched prefix is removed from the key
  inside the bucket.

## Profile accessor

```go
result, err := envsplit.SplitProfile(store, "production", rules, opts)
```

Loads a profile from any `ProfileGetter` and then splits it.
