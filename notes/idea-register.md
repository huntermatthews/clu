# Source Registration Refactoring Plan

## Goal

Refactor `internal/facts/sources/*.go` to be organized into OS-specific and generic packages while enabling cross-platform mock testing.

## Problem Statement

Currently all source files are in `internal/facts/sources/` and are referenced by OS-specific factory files (`internal/facts/linux.go`, `internal/facts/darwin.go`, `internal/facts/windows.go`).

**Current Issues:**
1. All sources in one directory makes organization unclear
2. OS factories have hand-curated lists of compatible sources
3. Build tags would break cross-platform mock testing

## Solution: Registration-Based Architecture

Use `init()` functions in each source file to self-register with a central registry, declaring which OS(es) they support. This solves the circular import problem and enables cross-platform mock testing.

## Proposed Directory Structure

``` text
internal/facts/
├── registry/
│   └── registry.go      # Central registration, no imports of source packages
├── types/
│   └── sources.go       # Interfaces only
├── generic/
│   └── clu.go           # Registers with: "linux", "darwin", "windows"
├── linux/
│   ├── dnf_checkupdate.go   # Registers with: "linux"
│   ├── ip_addr.go           # Registers with: "linux"
│   ├── ipmitool.go          # Registers with: "linux"
│   ├── lscpu.go             # Registers with: "linux"
│   ├── lsmem.go             # Registers with: "linux"
│   ├── no_salt.go           # Registers with: "linux"
│   ├── os_release.go        # Registers with: "linux"
│   ├── proc_cpuinfo.go      # Registers with: "linux"
│   ├── proc_cpuinfo2.go     # Registers with: "linux"
│   ├── proc_cpuinfo3.go     # Registers with: "linux"
│   ├── proc_uptime.go       # Registers with: "linux"
│   ├── selinux.go           # Registers with: "linux"
│   ├── sys_dmi.go           # Registers with: "linux"
│   ├── udevadm_ram.go       # Registers with: "linux"
│   └── virt_what.go         # Registers with: "linux"
├── unix/
│   ├── uname.go         # Registers with: "linux", "darwin", "freebsd"
│   └── uptime.go        # Registers with: "linux", "darwin", "freebsd"
├── macos/
│   ├── macos_name.go        # Registers with: "darwin"
│   ├── sw_vers.go           # Registers with: "darwin"
│   └── system_version_plist.go  # Registers with: "darwin"
└── windows/
    └── windows_systeminfo.go    # Registers with: "windows"
```

## Implementation Details

### 1. Create Registry Package

```go
// filepath: internal/facts/registry/registry.go
package registry

import (
    "github.com/huntermatthews/clu/internal/facts/types"
    "sync"
)

var (
    mu       sync.RWMutex
    sources  = make(map[string][]types.SourceFactory)
)

// SourceFactory is a function that creates a new Source instance
type SourceFactory func() types.Sources

// Register adds a source factory for the given operating systems
// Example: Register(NewLscpuFactory, "linux")
func Register(factory SourceFactory, oses ...string) {
    mu.Lock()
    defer mu.Unlock()

    for _, os := range oses {
        sources[os] = append(sources[os], factory)
    }
}

// GetSources returns all registered sources for the given OS
func GetSources(os string) []types.Sources {
    mu.RLock()
    defer mu.RUnlock()

    factories := sources[os]
    result := make([]types.Sources, len(factories))
    for i, factory := range factories {
        result[i] = factory()
    }
    return result
}

// Clear is useful for testing
func Clear() {
    mu.Lock()
    defer mu.Unlock()
    sources = make(map[string][]types.SourceFactory)
}
```

### 2. Update Each Source File

Each source file adds an `init()` function that registers itself:

#### Linux-only source

```go
// filepath: internal/facts/linux/lscpu.go
package linux

import (
    "github.com/huntermatthews/clu/internal/facts/registry"
    "github.com/huntermatthews/clu/internal/facts/types"
)

func init() {
    registry.Register(func() types.Sources {
        return &Lscpu{}
    }, "linux")
}

type Lscpu struct {
    // ...existing code...
}
```

#### Unix-compatible source

```go
// filepath: internal/facts/unix/uname.go
package unix

import (
    "github.com/huntermatthews/clu/internal/facts/registry"
    "github.com/huntermatthews/clu/internal/facts/types"
)

func init() {
    registry.Register(func() types.Sources {
        return &Uname{}
    }, "linux", "darwin", "freebsd")
}

type Uname struct {
    // ...existing code...
}
```

#### Cross-platform source

```go
// filepath: internal/facts/generic/clu.go
package generic

import (
    "github.com/huntermatthews/clu/internal/facts/registry"
    "github.com/huntermatthews/clu/internal/facts/types"
)

func init() {
    registry.Register(func() types.Sources {
        return &Clu{}
    }, "linux", "darwin", "windows", "freebsd")
}

type Clu struct {
    // ...existing code...
}
```

### 3. Update OpSys Factory

```go
// filepath: internal/facts/opsys.go
package facts

import (
    "github.com/huntermatthews/clu/internal/facts/registry"
    "github.com/huntermatthews/clu/internal/facts/types"

    // Import all source packages to trigger init() registration
    // These imports have no circular dependency because they only
    // call registry.Register(), they don't import internal/facts
    _ "github.com/huntermatthews/clu/internal/facts/generic"
    _ "github.com/huntermatthews/clu/internal/facts/linux"
    _ "github.com/huntermatthews/clu/internal/facts/macos"
    _ "github.com/huntermatthews/clu/internal/facts/unix"
    _ "github.com/huntermatthews/clu/internal/facts/windows"
)

func OpSysFactory() *OpSys {
    // ...existing OS detection...

    var defaults, early []string
    var osName string

    switch kernel {
    case "Darwin":
        osName = "darwin"
        defaults = darwinDefaults
        early = darwinEarly
    case "Linux":
        osName = "linux"
        defaults = linuxDefaults
        early = linuxEarly
    case "Windows":
        osName = "windows"
        defaults = windowsDefaults
        early = windowsEarly
    default:
        panic("unsupported operating system; got " + kernel)
    }

    sources := registry.GetSources(osName)
    return NewOpSys(sources, defaults, early)
}
```

## Import Dependency Flow (No Cycles!)

```text
internal/facts/opsys.go
    ↓ imports (blank)
internal/facts/linux/*.go
    ↓ imports
internal/facts/registry/registry.go
    ↓ imports
internal/facts/types/sources.go
```

**Key:** `internal/facts/linux/` never imports `internal/facts/`, only `registry` and `types`

## Benefits

### 1. Cross-Platform Mock Testing

- ✅ All packages compile on all platforms (no build tags)
- ✅ Linux sources can be tested on macOS with mock data
- ✅ Works with existing mock mode infrastructure in `cmd/clu/main.go`

### 2. Self-Documenting

- Each source file declares its own OS compatibility
- Easy to see which OSes a source supports by reading its `init()`

### 3. No Hand-Curated Lists

- OS factories are auto-generated from registrations
- Adding a new source just requires creating a file with `init()`

### 4. Clean Organization

- Sources grouped by OS compatibility
- Clear separation between generic, unix, and OS-specific sources

### 5. Testable

```go
// filepath: internal/facts/linux/lscpu_test.go
package linux_test

import (
    "testing"

    "github.com/huntermatthews/clu/internal/facts/linux"
    "github.com/huntermatthews/clu/internal/input"
)

func TestLscpuMockMode(t *testing.T) {
    // This test runs on ANY OS because no build tags!
    input.CommandRunner = input.MockTextProgram
    // ... setup mock data ...

    source := &linux.Lscpu{}
    err := source.Gather()
    // ... assertions ...
}
```

## Migration Steps

1. **Create new packages:**
   - `internal/facts/registry/` with `registry.go`
   - `internal/facts/types/` (move `Sources` interface from `internal/facts/sources/`)
   - `internal/facts/generic/`, `internal/facts/linux/`, `internal/facts/macos/`, `internal/facts/unix/`, `internal/facts/windows/`

2. **Move and update source files:**
   - Move each `.go` file to appropriate package directory
   - Add `init()` function with `registry.Register()` call
   - Update package name
   - Update imports to use `internal/facts/types` and `internal/facts/registry`

3. **Update factory files:**
   - Modify `internal/facts/opsys.go` to use blank imports and `registry.GetSources()`
   - Remove old factory files: `internal/facts/linux.go`, `internal/facts/darwin.go`, `internal/facts/windows.go`

4. **Update tests:**
   - Move test files to match new package structure
   - Update imports

5. **Clean up:**
   - Remove old `internal/facts/sources/` directory
   - Update documentation

## Future Extensibility

The registration system can be extended with metadata:

```go
type SourceRegistration struct {
    Factory  SourceFactory
    OSes     []string
    Priority int          // For ordering
    Tags     []string     // e.g., "requires-root", "slow", "network"
}

func RegisterWithMetadata(factory SourceFactory, metadata SourceRegistration)
```

This would enable:
- Filtering sources by tags in mock mode
- Skipping sources that require root in tests
- Ordering sources by priority
- Marking sources that require network access

## OS Compatibility Mapping

### Generic (all platforms)

- `clu.go` → `generic/` → registers: "linux", "darwin", "windows", "freebsd"

### Linux-only

- `dnf_checkupdate.go` → `linux/`
- `ip_addr.go` → `linux/`
- `ipmitool.go` → `linux/`
- `lscpu.go` → `linux/`
- `lsmem.go` → `linux/`
- `no_salt.go` → `linux/`
- `os_release.go` → `linux/`
- `proc_cpuinfo.go` → `linux/`
- `proc_cpuinfo2.go` → `linux/`
- `proc_cpuinfo3.go` → `linux/`
- `proc_uptime.go` → `linux/`
- `selinux.go` → `linux/`
- `sys_dmi.go` → `linux/`
- `udevadm_ram.go` → `linux/`
- `virt_what.go` → `linux/`

### Unix (Linux, macOS, FreeBSD)

- `uname.go` → `unix/` → registers: "linux", "darwin", "freebsd"
- `uptime.go` → `unix/` → registers: "linux", "darwin", "freebsd"

### macOS-only

- `macos_name.go` → `macos/`
- `sw_vers.go` → `macos/`
- `system_version_plist.go` → `macos/`

### Windows-only

- `windows_systeminfo.go` → `windows/`
