package main

// Go port of src/clu/cli.py using pflag for argument parsing.
// This provides a main-style entrypoint plus subcommand registration similar
// to the Python implementation. Subcommand implementations (report, archive,
// requires) are expected to be added separately.

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/subcmd"
)

// CLI defines the root command and global flags.
type CLI struct {
	Debug     bool                `help:"Enable debug logging."`
	Net       bool                `name:"net" help:"Enable network access."`
	Facts     subcmd.FactsCmd     `cmd:"" help:"Show facts (stub)."`
	Collector subcmd.CollectorCmd `cmd:"" help:"Run collector (stub)."`
	Requires  subcmd.RequiresCmd  `cmd:"" help:"Requires actions: list or check."`
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
		pkg.CluConfig.Debug = true
	}

	if cli.Net {
		fmt.Fprintln(os.Stderr, "net: enabled")
		pkg.CluConfig.NetEnabled = true
	}

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
