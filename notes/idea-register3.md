# Clu Registration and Extension Architecture — Phase One

## Goal

Refactor Clu so public and internal extensions use the same registration model:

- move the current `internal/` packages to `pkg/` so an internal module can
  import Clu's supported APIs;
- organize fact sources by OS compatibility;
- have every active public fact source self-register from `init()`;
- have every built-in and internal subcommand self-register from `init()`;
- retain cross-platform mock testing and complete `collector` fixtures; and
- let the internal Clu branch include private sources and commands through one
  private aggregate import.

There is one `clu` binary. The public branch has public registrations only; the
internal branch adds a private module dependency and imports its registrations.

## Public Package Layout

```text
pkg/
├── facts/
│   ├── opsys.go
│   ├── registry/
│   │   └── registry.go
│   ├── types/
│   │   ├── facts.go
│   │   ├── provides.go
│   │   ├── requires.go
│   │   └── sources.go
│   ├── generic/
│   ├── linux/
│   ├── macos/
│   ├── unix/
│   └── windows/
├── global/
├── input/
└── subcmd/
    ├── registry.go
    ├── facts.go
    ├── collector.go
    ├── requires.go
    ├── tools.go
    ├── version.go
    ├── check.go
    └── help.go

cmd/clu/
└── main.go
```

The package move is mechanical:

```text
internal/facts/   -> pkg/facts/
internal/input/   -> pkg/input/
internal/subcmd/  -> pkg/subcmd/
internal/global/  -> pkg/global/
```

`pkg/` is a convention, not a visibility boundary. These are deliberate APIs
for the internal extension module.

## Fact Source Registration

### Source contract

`pkg/facts/types` retains the current source contract:

```go
type Sources interface {
    Provides(Provides)
    Requires(*Requires)
    Parse(*FactDB)
}
```

Sources use `pkg/input.CommandRunner` and `pkg/input.FileReader` rather than
direct OS calls. This preserves the current `collector` and `--mock-dir`
behavior for both public and internal sources.

### Minimal registry

`pkg/facts/registry` is a minimal OS-to-factory registry:

```go
package registry

import (
    "sync"

    "github.com/huntermatthews/clu/pkg/facts/types"
)

type SourceFactory func() types.Sources

var (
    mu      sync.RWMutex
    sources = make(map[string][]SourceFactory)
)

func Register(factory SourceFactory, oses ...string) {
    mu.Lock()
    defer mu.Unlock()

    for _, osName := range oses {
        sources[osName] = append(sources[osName], factory)
    }
}

func GetSources(osName string) []types.Sources {
    mu.RLock()
    factories := append([]SourceFactory(nil), sources[osName]...)
    mu.RUnlock()

    result := make([]types.Sources, 0, len(factories))
    for _, factory := range factories {
        result = append(result, factory())
    }
    return result
}
```

Phase one deliberately has no source IDs, priorities, metadata, collision
validation, sorting, or global `Clear()` function. Those enhancements are
reserved for `idea-register-phase2.md`.

### Init-based self-registration

Every active **public source** registers itself in its own package/file. The
source declares its supported OSes locally:

```go
// pkg/facts/linux/lscpu.go
package linux

import (
    "github.com/huntermatthews/clu/pkg/facts/registry"
    "github.com/huntermatthews/clu/pkg/facts/types"
)

func init() {
    registry.Register(func() types.Sources {
        return &Lscpu{}
    }, "linux")
}
```

A shared source can register for more than one OS:

```go
// pkg/facts/generic/clu.go
func init() {
    registry.Register(func() types.Sources {
        return &Clu{}
    }, "linux", "darwin", "windows")
}
```

This is the key mechanism that removes hand-curated source lists. Adding an
active source requires placing it in the appropriate package and adding its
local `init()` registration.

### Source organization and wrappers

Only actual fact providers self-register. Private collection strategies do not.

| Registered public wrapper | OS | Fact area | Internal collector order |
| --- | --- | --- | --- |
| `Memory` | Linux | `phy.ram` | `lsmem` -> `udevadm` -> future `ProcMeminfo` |
| `Platform` | Linux | `phy.platform` | `virt-what` -> stubbed `systemd-detect-virt` |
| `SystemVersion` | Darwin | OS version facts | SystemVersion.plist -> `sw_vers` |

The private collector types (`lsmem`, `udevadmRam`, `virtWhat`,
`systemdDetectVirt`, `systemVersionPlist`, and `swVers`) are implementation
details. They do not implement `Provides` or `Requires` and never self-register.

Fallback rules:

- a collector succeeds only with a usable non-empty, non-`Unknown/Error` value;
- a later successful collector may replace an earlier `Unknown/Error` value;
- `FactDB.AddFact` must add a key to its tier index only once when replacing a
  value;
- each wrapper includes requirements for every collector, not merely the first
  one, because `collector` must capture all possible source inputs;
- `ProcMeminfo`, not `ProcCpuinfo`, is the future final RAM fallback;
- the `systemd-detect-virt` parser is a current stub: its command is collected,
  but it does not yet provide a fallback value; and
- `sw_vers` cannot provide `id.build_id`, so that fact remains `Unknown/Error`
  if it is used after plist parsing fails.

`macos_checkupdate.go` is currently inactive. It must not self-register until
it is deliberately adopted as an active Darwin source.

### OpSys bootstrap

`OpSysFactory` still determines the OS and supplies early-fact settings. It
loads built-in source packages solely to execute their registrations, then
asks the registry for the detected OS:

```go
import (
    "github.com/huntermatthews/clu/pkg/facts/registry"
    "github.com/huntermatthews/clu/pkg/facts/types"
    factunix "github.com/huntermatthews/clu/pkg/facts/unix"

    _ "github.com/huntermatthews/clu/pkg/facts/generic"
    _ "github.com/huntermatthews/clu/pkg/facts/linux"
    _ "github.com/huntermatthews/clu/pkg/facts/macos"
    _ "github.com/huntermatthews/clu/pkg/facts/windows"
)
```

`unix` is a normal import because `OpSysFactory` uses `&factunix.Uname{}` to
detect a non-Windows kernel before it knows whether the OS is Linux or Darwin.
The remaining imports are blank imports that trigger `init()` registration.

Source packages import `registry`, `types`, `input`, and other lower-level
packages, but never the parent `facts` package. This avoids import cycles.

## Dynamic Subcommand Registration

### Root CLI

All commands, including built-in commands, use dynamic registration. The root
`CLI` struct has global flags only:

```go
type CLI struct {
    Debug   bool             `help:"Enable debug logging."`
    Net     bool             `name:"net" help:"Enable network access."`
    MockDir string           `help:"Enable mock mode for testing." hidden:""`
    Version kong.VersionFlag `help:"Print version information and quit."`
}
```

The static command fields for `FactsCmd`, `CollectorCmd`, `RequiresCmd`,
`ToolsCmd`, `VersionCmd`, `CheckCmd`, and `HelpCmd` are removed.

### Minimal command registry

`pkg/subcmd` holds the command registry:

```go
type Registration struct {
    Name    string
    Help    string
    Command any
    Tags    []string
}

func Register(registration Registration)
func KongOptions() []kong.Option
```

`KongOptions()` creates one `kong.DynamicCommand()` option per registration.
Phase one has no command IDs, groups, sorting, metadata, or custom duplicate
validation.

Every built-in command self-registers from its own file:

```go
func init() {
    Register(Registration{
        Name:    "facts",
        Help:    "Show facts.",
        Command: &FactsCmd{},
        Tags:    []string{`default:"withargs"`},
    })
}
```

The `default:"withargs"` tag preserves the current behavior of running
`facts` for a bare `clu` invocation. The global `--version` flag remains
separate from the dynamically registered `version` command.

Built-in registrations are:

| Type | Name |
| --- | --- |
| `FactsCmd` | `facts` |
| `CollectorCmd` | `collector` |
| `RequiresCmd` | `requires` |
| `ToolsCmd` | `tools` |
| `VersionCmd` | `version` |
| `CheckCmd` | `check` |
| `HelpCmd` | `help` |

`cmd/clu/main.go` builds Kong from the registered command options:

```go
options := []kong.Option{
    kong.Name("clu"),
    kong.UsageOnError(),
    kong.Writers(stdout, stderr),
    kong.BindTo(stdout, (*input.Stdout)(nil)),
    kong.BindTo(stderr, (*input.Stderr)(nil)),
}
options = append(options, subcmd.KongOptions()...)

k, err := kong.New(&CLI{}, options...)
```

## Internal Extension Repository

The internal Clu branch adds the private module requirement to `go.mod` and
imports one aggregate registration package:

```go
import _ "git.example.com/platform/clu-private/register"
```

The aggregate package blank-imports private source and subcommand packages.
Their `init()` functions register them before `main()` calls `run()`:

```go
package register

import (
    _ "git.example.com/platform/clu-private/sources/inventory"
    _ "git.example.com/platform/clu-private/subcmd/report"
)
```

A private fact source imports `pkg/facts/registry`, `pkg/facts/types`, and
`pkg/input`, then self-registers exactly like a built-in source. A private
subcommand imports `pkg/subcmd` and self-registers exactly like a built-in
command.

The private module must be pinned to a known version and fetched through the
internal Git configuration (`GOPRIVATE` as appropriate).

## Testing Requirements

- Source implementations must not use Go OS build tags or `_linux.go`,
  `_darwin.go`, or `_windows.go` filename suffixes; this preserves
  cross-platform mock tests.
- Test that all built-in sources and commands self-register when their packages
  are imported.
- Test that a bare `clu` invocation still runs `facts`.
- Test the internal branch with its aggregate registration import, private
  source requirements, collector output, and `--mock-dir` fixtures.
- Test fallback wrappers with a successful preferred collector, successful
  fallback collector, and all-collector failure.

## Migration Order

1. Move `internal/facts`, `internal/input`, `internal/subcmd`, and
   `internal/global` to their corresponding `pkg/` locations; update imports.
2. Create `pkg/facts/registry` with the minimal `Register` and `GetSources`
   API.
3. Move built-in source files into generic, Linux, Unix, macOS, and Windows
   packages. Add an `init()` registration to each active public source only.
4. Update `OpSysFactory` to blank-import built-in source packages and construct
   its sources with `registry.GetSources(osName)`, retaining OS detection and
   early-fact behavior.
5. Create the minimal `pkg/subcmd` dynamic registry.
6. Add an `init()` registration for every built-in command; update
   `cmd/clu/main.go` to create Kong from `subcmd.KongOptions()`.
7. In the internal Clu branch, add the private module requirement and aggregate
   blank import.
8. Run public and internal integration tests, including `collector` fixture
   generation and `--mock-dir` execution.

## Deferred Work

`idea-register-phase2.md` covers optional later features: source/command IDs,
priorities, metadata/tags, explicit overrides, duplicate validation,
deterministic ordering, and isolated registry test helpers.
