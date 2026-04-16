# inject

The `inject` package runs subprocesses with a profile's environment variables
merged on top of the current process environment.

## Usage

```go
inj := inject.NewInjector(os.Environ())
cmd := inj.Command(profile.Vars, "make", "build")
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Run()
```

## Behaviour

- The base environment (typically `os.Environ()`) is used as the starting point.
- Keys present in the overlay (profile vars) **override** base values.
- Keys present only in the base are preserved unchanged.
- The returned `*exec.Cmd` is ready to run; callers wire up `Stdin`/`Stdout`/`Stderr` as needed.
