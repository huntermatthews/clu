# cliv1 example

This example uses `urfave/cli` v1.

## Run

```bash
# Help
go run ./examples/cliv1 --help
go run ./examples/cliv1 facts --help

# Facts: tier 2 with network enabled
go run ./examples/cliv1 facts -t 2 --net

# Requires: list or check
go run ./examples/cliv1 requires list
go run ./examples/cliv1 requires check
```

## Run tests

```bash
# From repo root
go test ./examples/cliv1 -run .

# Or from this folder
cd examples/cliv1
go test -run .
```

## urfave/cli v1 Cheat Sheet

- App & commands:
 	- `cli.NewApp()` — <https://pkg.go.dev/github.com/urfave/cli#NewApp>
 	- `cli.App` — <https://pkg.go.dev/github.com/urfave/cli#App>
 	- `cli.Command` — <https://pkg.go.dev/github.com/urfave/cli#Command>
- Flags:
 	- `cli.IntFlag` — <https://pkg.go.dev/github.com/urfave/cli#IntFlag>
 	- `cli.BoolFlag` — <https://pkg.go.dev/github.com/urfave/cli#BoolFlag>
 	- `cli.StringFlag` — <https://pkg.go.dev/github.com/urfave/cli#StringFlag>
- Lifecycle & run:
 	- `app.Before` — <https://pkg.go.dev/github.com/urfave/cli#BeforeFunc>
 	- `command.Action` — <https://pkg.go.dev/github.com/urfave/cli#ActionFunc>
 	- `app.Run(os.Args)` — <https://pkg.go.dev/github.com/urfave/cli#App.Run>
