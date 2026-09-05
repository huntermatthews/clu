# Clu Registration Architecture — Phase Two

## Purpose

This document contains the registration features intentionally deferred from
phase one in `idea-register.md`. Implement them only after the minimal
public/private source and command registration model is working in both the
public and internal Clu branches.

Phase two addresses intentional source overrides, diagnostic quality,
deterministic ordering, command grouping, and richer registration metadata.

## Source Registry Enhancements

### Source IDs

Replace the phase-one factory-only registration API with a registration value:

```go
type SourceFactory func() types.Sources

type Registration struct {
    ID       string
    Factory  SourceFactory
    OSes     []string
    Priority int
    Tags     []string
}

func Register(reg Registration)
```

`ID` is a stable unique source identifier, such as:

```text
clu.linux.memory
clu.darwin.system-version
company.asset-inventory
```

It improves startup diagnostics and lets the registry reject accidental
duplicate registrations.

### Priority and deterministic ordering

Phase one requires public source providers to use distinct fact keys. Phase
two may support an intentional source override, such as an internal source
replacing a public provider for a fact key.

The registry sorts registrations by ascending priority, then by ID:

```go
sort.Slice(registrations, func(i, j int) bool {
    if registrations[i].Priority != registrations[j].Priority {
        return registrations[i].Priority < registrations[j].Priority
    }
    return registrations[i].ID < registrations[j].ID
})
```

This avoids relying on import or `init()` order. Source precedence remains
visible at the registration site.

The exact overwrite direction must be documented and tested. With the current
`OpSys.Provides()` implementation, the provider processed later replaces an
earlier source for the same fact key. Therefore, under ascending priority, a
higher numeric priority wins if later registrations overwrite earlier ones.

### Duplicate and override policy

Phase two should validate registrations at startup:

- duplicate registration IDs are errors;
- duplicate fact keys for one OS are errors by default;
- an explicit override tag permits a duplicate fact key;
- the error message identifies the OS, fact key, and both source IDs.

For example:

```go
registry.Register(registry.Registration{
    ID:       "company.asset-owner",
    Factory:  func() types.Sources { return &AssetOwner{} },
    OSes:     []string{"linux"},
    Priority: 500,
    Tags:     []string{"override"},
})
```

Private facts should still normally use a namespace such as `company.*`.
Overrides should be exceptional and documented alongside their priority.

### Metadata

`Tags` are optional source metadata. They can later identify collection cost or
constraints without adding special-case APIs:

```text
requires-root
slow
network
internal
```

Possible future uses:

- list sources and their requirements for diagnostics;
- filter slow or networked sources in tests;
- skip root-only collectors when requested;
- report whether a fact came from an internal extension;
- apply policy to private-source overrides.

Do not use tags to recreate fallback behavior. Fallbacks remain local to a
single public wrapper source (`Memory`, `Platform`, or `SystemVersion`).

### Test isolation

Do not expose a global `Clear()` for an init-populated registry. `init()` runs
once per test process, so a reset cannot recreate normal registrations.

If phase two needs isolated registry tests, use one of these approaches:

1. make the registry an instantiable type and test fresh instances;
2. use an unexported test helper that snapshots and restores all registration
   state; or
3. use explicit registration functions in test packages rather than global
   `init()` registration.

## Command Registry Enhancements

### Command IDs and groups

Replace the phase-one command registration value with:

```go
type Registration struct {
    ID      string
    Name    string
    Help    string
    Group   string
    Command any
    Tags    []string
}
```

`ID` gives a stable diagnostic identity independent of the user-facing command
name. `Group` can organize Kong help output, for example:

```text
Built-in commands
Internal commands
```

### Duplicate validation

Before passing registrations to `kong.DynamicCommand()`, validate:

- every ID is unique;
- every root command name is unique;
- registration has a non-nil command value.

A private command must not silently replace a built-in command such as
`facts`, `collector`, or `check`.

### Deterministic command listing

Sort registrations by command name, or by group then name, before creating
Kong options. This keeps generated help and related tests stable without
relying on `init()` order.

## Phase-Two Completion Criteria

Phase two is complete when:

1. source and command registrations have stable IDs;
2. source order is explicit and deterministic;
3. accidental duplicate registrations fail with actionable diagnostics;
4. intentional source overrides require an explicit documented opt-in;
5. command help is deterministic and can be grouped;
6. registry tests use isolated state without a public destructive reset; and
7. public and internal branch integration tests cover registration conflicts
   and permitted overrides.
