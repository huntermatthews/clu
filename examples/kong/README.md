# Kong CLI Example

A self-contained example using [alecthomas/kong] to wire a multi-subcommand CLI.
This example does not reference other code in the repository.

## Run

```sh
# Help
go run ./examples/kong --help
go run ./examples/kong facts --help

# Facts: tier 2 with network enabled
go run ./examples/kong facts -t 2 --net

# Output format
go run ./examples/kong facts --out json

# Requires: list or check
go run ./examples/kong requires list
go run ./examples/kong requires check

# Debug flag (prints to stderr)
go run ./examples/kong --debug facts -t 1
```

### Environment variables

You can use `kong.DefaultEnvars("CLU")` to allow flags to be set via environment variables. For example, adding this to `kong.New(...)` enables:

```sh
export CLU_TIER=3
export CLU_NET=true
go run ./examples/kong facts
# Interprets env vars as defaults, equivalent to: facts --tier 3 --net
```

## Commands

- Root flags: `--help` (from Kong), `--debug`
- Subcommands:
  - `facts` with `-t/--tier` (1,2,3), `--out` (dots,json,shell), `--net`
  - `collector` (stub)
  - `requires <list|check>` (required positional)

## Quick start

Ensure the dependency is available:

```sh
go get github.com/alecthomas/kong@latest
```

Run the example:

```sh
# Help
go run ./examples/kong --help

# Facts: tier 2 with network enabled
go run ./examples/kong facts -t 2 --net

# Collector stub
go run ./examples/kong collector

# Requires stub
go run ./examples/kong requires list

# Debug flag (prints to stderr)
go run ./examples/kong --debug facts -t 1
```

## Kong Options Cheat Sheet

You can pass these option constructors to `kong.New(...)` to configure behavior:

- Name & description:
  - `kong.Name(string)` — https://pkg.go.dev/github.com/alecthomas/kong#Name
  - `kong.Description(string)` — <https://pkg.go.dev/github.com/alecthomas/kong#Description>
  - `kong.Summary(string)` — https://pkg.go.dev/github.com/alecthomas/kong#Summary
  - `kong.Epilogue(string)` — https://pkg.go.dev/github.com/alecthomas/kong#Epilogue
- Help & usage:
  - `kong.UsageOnError()` — https://pkg.go.dev/github.com/alecthomas/kong#UsageOnError
  - `kong.NoDefaultHelp()` — https://pkg.go.dev/github.com/alecthomas/kong#NoDefaultHelp
  - `kong.Help(io.Writer)` — https://pkg.go.dev/github.com/alecthomas/kong#Help
  - `kong.ConfigureHelp(kong.HelpOptions)` — https://pkg.go.dev/github.com/alecthomas/kong#ConfigureHelp
  - `kong.HelpFlag(kong.HelpFlag)` — https://pkg.go.dev/github.com/alecthomas/kong#HelpFlag
  - `kong.Writer(io.Writer)` — https://pkg.go.dev/github.com/alecthomas/kong#Writer
- Env & variables:
  - `kong.DefaultEnvars(prefix string)` — https://pkg.go.dev/github.com/alecthomas/kong#DefaultEnvars
  - `kong.Vars(map[string]string)` — https://pkg.go.dev/github.com/alecthomas/kong#Vars
- Flags & parsing:
  - `kong.NoShorts()` — https://pkg.go.dev/github.com/alecthomas/kong#NoShorts
  - `kong.ShortFlagDelimiter(rune)` — https://pkg.go.dev/github.com/alecthomas/kong#ShortFlagDelimiter
  - `kong.Model(*kong.Model)` — https://pkg.go.dev/github.com/alecthomas/kong#Model
  - `kong.Resolvers(kong.Resolver...)` — https://pkg.go.dev/github.com/alecthomas/kong#Resolvers
- Execution & exit:
  - `kong.Bind(any)` — https://pkg.go.dev/github.com/alecthomas/kong#Bind
  - `kong.Exit(func(code int))` — https://pkg.go.dev/github.com/alecthomas/kong#Exit
  - `kong.Usage(func(*kong.Context))` — https://pkg.go.dev/github.com/alecthomas/kong#Usage


## Tests

```sh
# Run the tiny test suite for this example
go test ./examples/kong -v
```

[alecthomas/kong]: https://github.com/alecthomas/kong
