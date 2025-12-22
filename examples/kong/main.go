package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// CLI defines the root command and global flags.
type CLI struct {
	Debug     bool         `help:"Enable debug logging."`
	Facts     FactsCmd     `cmd:"" help:"Show facts (stub)."`
	Collector CollectorCmd `cmd:"" help:"Run collector (stub)."`
	Requires  RequiresCmd  `cmd:"" help:"Requires actions: list or check."`
}

// FactsCmd implements the "facts" subcommand.
type FactsCmd struct {
	Tier         int      `name:"tier" short:"t" enum:"1,2,3" default:"1" help:"Tier level (1, 2, or 3)."`
	OutputFormat string   `name:"out" enum:"dots,json,shell" default:"dots" help:"Output format: dots, json, or shell."`
	Net          bool     `name:"net" help:"Enable network access."`
	Facts        []string `arg:"" optional:"true" help:"Zero or more fact names."`
}

func (f *FactsCmd) Run() error {
	tier := f.Tier
	if len(f.Facts) > 0 {
		fmt.Printf("facts: tier=%d net=%v (stub) %s\n", tier, f.Net, strings.Join(f.Facts, " "))
	} else {
		fmt.Printf("facts: tier=%d net=%v (stub)\n", tier, f.Net)
	}
	return nil
}

// CollectorCmd implements the "collector" subcommand (stub only).
type CollectorCmd struct {
	OutputDir string `name:"output-dir" default:"/tmp" help:"Directory to write the archive/output."`
}

func (c *CollectorCmd) Run() error {
	// Stub action; show chosen output directory to confirm flag wiring.
	fmt.Printf("collector: running (stub) output-dir=%s\n", c.OutputDir)
	return nil
}

// RequiresCmd implements the "requires" subcommand (stub only).
type RequiresCmd struct {
	Mode string `arg:"" enum:"list,check" help:"Operation to perform: list or check."`
}

func (r *RequiresCmd) Run() error {
	switch r.Mode {
	case "list":
		fmt.Println("requires: listing (stub)")
	case "check":
		fmt.Println("requires: checking (stub)")
	}
	return nil
}

func main() {
	cli := &CLI{}
	k, err := kong.New(cli,
		kong.Name("clu"),
		kong.Description("Kong example with facts/collector/requires subcommands."),
		kong.DefaultEnvars("CLU"),
		kong.UsageOnError(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init CLI:", err)
		os.Exit(2)
	}

	ctx, err := k.Parse(os.Args[1:])
	if err != nil {
		k.FatalIfErrorf(err)
	}

	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: enabled")
	}

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}

// Note: Kong does not support numeric-only short options like -1/-2/-3.
// Use -t/--tier with values 1,2,3 instead (e.g., -t=2 or --tier 2).
