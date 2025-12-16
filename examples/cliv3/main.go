package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	cli "github.com/urfave/cli/v3"
)

// Note: urfave/cli does not support numeric-only short flags for options.
// Use named flags like --one/--two/--three or a single flag with a value.

func main() {
	root := &cli.Command{
		Name:  "clu",
		Usage: "urfave/cli v3 example with facts/collector/requires",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug", Usage: "Enable debug logging."},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.Bool("debug") {
				// Match the Kong example behavior
				fmt.Fprintln(os.Stderr, "debug: enabled")
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "facts",
				Usage: "Show facts (stub)",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "tier", Aliases: []string{"t"}, Usage: "Tier level (1, 2, or 3).", Value: 1},
					&cli.BoolFlag{Name: "net", Usage: "Enable network access."},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					tier := cmd.Int("tier")
					if tier < 1 || tier > 3 {
						return fmt.Errorf("invalid --tier %d (allowed: 1, 2, 3)", tier)
					}
					names := cmd.Args().Slice()
					if len(names) > 0 {
						fmt.Printf("facts: tier=%d net=%v (stub) %s\n", tier, cmd.Bool("net"), strings.Join(names, " "))
					} else {
						fmt.Printf("facts: tier=%d net=%v (stub)\n", tier, cmd.Bool("net"))
					}
					return nil
				},
			},
			{
				Name:  "collector",
				Usage: "Run collector (stub)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("collector: running (stub)")
					return nil
				},
			},
			{
				Name:  "requires",
				Usage: "Requires actions: list or check",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					if len(args) < 1 {
						return fmt.Errorf("requires expects one positional argument: list or check")
					}
					mode := args[0]
					switch mode {
					case "list":
						fmt.Println("requires: listing (stub)")
					case "check":
						fmt.Println("requires: checking (stub)")
					default:
						return fmt.Errorf("invalid requires mode %q (allowed: list, check)", mode)
					}
					return nil
				},
			},
		},
	}

	if err := root.Run(context.Background(), os.Args); err != nil {
		// Use the library's default error printing convention.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
