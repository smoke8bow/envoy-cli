# envclip

Provides clipboard-style copy/paste for environment variable subsets.

## Overview

`envclip` lets you copy a selection of keys from one profile's variable map
into an in-process clipboard, then paste them into any other map — useful when
building CLI commands that need to transfer a subset of variables between
profiles without a full merge.

## Usage

```go
m := envclip.NewManager()

// Copy all vars from the "dev" profile map.
_ = m.Copy("dev", devVars, nil)

// Copy only selected keys.
_ = m.Copy("dev", devVars, []string{"DB_HOST", "DB_PORT"})

// Paste into another map (does not mutate the destination).
result, err := m.Paste(stagingVars)

// Inspect what is on the clipboard.
keys, _ := m.Keys()

// Clear when done.
m.Clear()
```

## Behaviour

- `Copy` with a `nil` or empty keys slice copies **all** keys.
- `Paste` **overwrites** matching keys in the destination; keys present only in
  the destination are preserved.
- The original destination map is **never mutated**.
- `Keys` returns keys in sorted order.
