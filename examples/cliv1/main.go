package main

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "clu"
	app.Usage = "urfave/cli v1 example with facts/collector/requires"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "Enable debug logging."},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			fmt.Fprintln(os.Stderr, "debug: enabled")
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "facts",
			Usage: "Show facts (stub)",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "tier, t", Usage: "Tier level (1, 2, or 3).", Value: 1},
				cli.StringFlag{Name: "out", Usage: "Output format: dots, json, or shell.", Value: "dots"},
				cli.BoolFlag{Name: "net", Usage: "Enable network access."},
			},
			Action: func(c *cli.Context) error {
				tier := c.Int("tier")
				if tier < 1 || tier > 3 {
					return fmt.Errorf("invalid --tier %d (allowed: 1, 2, 3)", tier)
				}
				outfmt := c.String("out")
				switch outfmt {
				case "dots", "json", "shell":
					// ok
				default:
					return fmt.Errorf("invalid --out %q (allowed: dots, json, shell)", outfmt)
				}
				names := c.Args()
				if len(names) > 0 {
					fmt.Printf("facts: tier=%d net=%v (stub) %s\n", tier, c.Bool("net"), strings.Join(names, " "))
				} else {
					fmt.Printf("facts: tier=%d net=%v (stub)\n", tier, c.Bool("net"))
				}
				return nil
			},
		},
		{
			Name:  "collector",
			Usage: "Run collector (stub)",
			Action: func(c *cli.Context) error {
				fmt.Println("collector: running (stub)")
				return nil
			},
		},
		{
			Name:  "requires",
			Usage: "Requires actions: list or check",
			Action: func(c *cli.Context) error {
				args := c.Args()
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
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
