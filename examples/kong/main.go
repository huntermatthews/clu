package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

const version = "0.1.0"

// CLI defines the root command and global flags.
type CLI struct {
	Debug     bool             `help:"Enable debug logging."`
	Verbose   bool             `help:"Enable verbose output."`
	Version   kong.VersionFlag `help:"Print version information and quit."`
	Facts     FactsCmd         `cmd:"" help:"Show facts." default:"withargs"`
	Collector CollectorCmd     `cmd:"" help:"Run collector."`
	Requires  RequiresCmd      `cmd:"" help:"Requires actions: list or check."`
}

// FactsCmd implements the "facts" subcommand.
type FactsCmd struct {
	Tier         int      `name:"tier" short:"t" enum:"1,2,3" default:"1" help:"Tier level (1, 2, or 3)."`
	OutputFormat string   `name:"out" enum:"dots,json,shell" default:"dots" help:"Output format: dots, json, or shell."`
	Net          bool     `name:"net" help:"Enable network access."`
	Facts        []string `arg:"" optional:"true" help:"Zero or more fact names."`
}

func (f *FactsCmd) Run(cli *CLI) error {
	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: facts run")
	}
	if cli.Verbose {
		fmt.Fprintln(os.Stderr, "verbose: facts run")
	}
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

func (c *CollectorCmd) Run(cli *CLI) error {
	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: collector run")
	}
	if cli.Verbose {
		fmt.Fprintln(os.Stderr, "verbose: collector run")
	}
	fmt.Println("collector: running (stub)")
	return nil
}

// RequiresCmd implements the "requires" subcommand (stub only).
type RequiresCmd struct {
	Mode string `arg:"" enum:"list,check" help:"Operation to perform: list or check."`
}

func (r *RequiresCmd) Run(cli *CLI) error {
	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: requires run")
	}
	if cli.Verbose {
		fmt.Fprintln(os.Stderr, "verbose: requires run")
	}

	switch r.Mode {
	case "list":
		fmt.Println("requires: listing (stub)")
	case "check":
		fmt.Println("requires: checking (stub)")
	}
	return nil
}

func main() {
	slog.Debug("Starting clu application")
	slog.Info("Starting clu application")
	slog.Warn("Starting clu application")
	slog.Error("Starting clu application")

	cli := &CLI{}
	k, err := kong.New(cli,
		kong.Name("clu"),
		kong.Description("Kong example with facts/collector/requires subcommands."),
		kong.Bind(cli),
		kong.UsageOnError(),
		kong.Vars{"version": "v0.1.0"},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to init CLI:", err)
		os.Exit(2)
	}

	ctx, err := k.Parse(os.Args[1:])
	if err != nil {
		k.FatalIfErrorf(err)
	}

	// Handle --version at root: print and exit.
	if cli.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if cli.Debug {
		fmt.Fprintln(os.Stderr, "debug: enabled")
	}
	if cli.Verbose {
		fmt.Fprintln(os.Stderr, "verbose: enabled")
	}

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
