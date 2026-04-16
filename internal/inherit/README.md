# inherit

The `inherit` package allows a profile to inherit environment variables from a
parent profile. Keys already present in the child are **never** overwritten —
only missing keys are filled in from the parent.

## Usage

```go
inh := inherit.NewInheritor(store)

// Preview the merged result without persisting:
merged, err := inh.Apply("base", "production")

// Merge and save back to the child profile:
merged, err = inh.Commit("base", "production")
```

## Behaviour

| Key present in child | Key present in parent | Result in child |
|---|---|---|
| yes | yes | child value kept |
| no  | yes | parent value copied in |
| yes | no  | child value kept |
