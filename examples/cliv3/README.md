# cliv3 example

This example uses `urfave/cli/v3`.

## Run

```bash
# Help
go run ./examples/cliv3 --help
go run ./examples/cliv3 facts --help

# Facts: tier 2 with network enabled
go run ./examples/cliv3 facts --tier 2 --net

# Requires: list or check
go run ./examples/cliv3 requires list
go run ./examples/cliv3 requires check
```

## Run tests

```bash
# From repo root
go test ./examples/cliv3 -run .

# Or from this folder
cd examples/cliv3
go test -run .
```

## urfave/cli v3 Cheat Sheet

- Commands & flags:
 	- `cli.Command` — <https://pkg.go.dev/github.com/urfave/cli/v3#Command>
 	- `cli.IntFlag` — <https://pkg.go.dev/github.com/urfave/cli/v3#IntFlag>
 	- `cli.BoolFlag` — <https://pkg.go.dev/github.com/urfave/cli/v3#BoolFlag>
 	- `cli.StringFlag` — <https://pkg.go.dev/github.com/urfave/cli/v3#StringFlag>
- App lifecycle:
 	- `cli.Command.Before` — <https://pkg.go.dev/github.com/urfave/cli/v3#BeforeFunc>
 	- `cli.Command.Action` — <https://pkg.go.dev/github.com/urfave/cli/v3#ActionFunc>
 	- `(*cli.Command).Run(ctx, args)` — <https://pkg.go.dev/github.com/urfave/cli/v3#Command.Run>
- Context:
 	- `cmd.Bool("flag")/Int/String` — <https://pkg.go.dev/github.com/urfave/cli/v3#Context>
 	- `cmd.Args()` — <https://pkg.go.dev/github.com/urfave/cli/v3#Args>
